import { NextResponse } from "next/server";
import { stripe } from "@/lib/stripe";
import { getServerSession } from "next-auth";
import { authOptions } from "@/lib/auth";
import { prisma } from "@/lib/db";
export const runtime = 'nodejs';
export async function POST(){
  const session = await getServerSession(authOptions);
  if(!session?.user?.email) return NextResponse.json({ error: "unauth" }, { status: 401 });
  const user = await prisma.user.findUnique({ where: { email: session.user.email } });
  const sub = await prisma.subscription.findFirst({ where: { userId: user!.id, status: { in: ["active","trialing","past_due"] } } });
  if(!sub) return NextResponse.json({ error: "no-subscription" }, { status: 404 });
  const portal = await stripe.billingPortal.sessions.create({
    customer: (await stripe.subscriptions.retrieve(sub.stripeSubId)).customer as string,
    return_url: `${process.env.NEXT_PUBLIC_SITE_URL}/settings/billing`,
  });
  return NextResponse.json({ url: portal.url });
}
