import { NextResponse } from "next/server";
import { getServerSession } from "next-auth";
import { authOptions } from "@/lib/auth";
import { recordUsage, checkQuota } from "@/lib/usage";
export const runtime = 'nodejs';
export async function POST(req: Request){
  const session = await getServerSession(authOptions);
  if(!session) return NextResponse.json({ error: 'unauth' }, { status: 401 });
  const form = await req.formData();
  const prompt = String(form.get('prompt')||'').slice(0, 8000);
  const userId = (session as any).userId as string;
  const quota = await checkQuota(userId);
  if(!quota.ok) return NextResponse.json({ error: 'limit' }, { status: 429 });
  const answer = `Echo: ${prompt.substring(0,200)}…`; // TODO: integrate your LLM
  await recordUsage(userId, 'chat_request', prompt.length, answer.length, { model: 'placeholder' });
  return NextResponse.json({ answer });
}
