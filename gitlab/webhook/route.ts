import { NextResponse } from 'next/server';
import { prisma } from '@/lib/db';
export const runtime = 'nodejs';
export async function POST(req: Request){
  const token = req.headers.get('x-gitlab-token');
  if((process.env.GITLAB_WEBHOOK_TOKEN ?? '') !== token) return new NextResponse('Invalid token', { status: 401 });
  const evt = await req.json();
  if(evt.object_kind === 'merge_request' && ['open','reopen','update'].includes(evt.object_attributes?.action)){
    const repoUrl = evt.object_attributes?.url || evt.project?.web_url;
    const job = await prisma.reviewJob.create({ data: { repoUrl, status: 'queued', meta: { provider: 'gitlab', mr: evt.object_attributes?.iid, class: 'quick' } } });
    return NextResponse.json({ ok:true, job });
  }
  return NextResponse.json({ ok:true });
}
