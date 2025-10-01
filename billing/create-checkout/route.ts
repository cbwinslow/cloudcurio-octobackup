import { NextResponse } from "next/server";
import { stripe, STRIPE_PRICE_PRO } from "@/lib/stripe";
import { getServerSession } from "next-auth";
import { authOptions } from "@/lib/auth";
import { prisma } from "@/lib/db";
export const runtime = 'nodejs';
export async function POST(){
  const session = await getServerSession(authOptions);
  if(!session?.user?.email) return NextResponse.json({ error: "unauth" }, { status: 401 });
  const user = await prisma.user.findUnique({ where: { email: session.user.email } });
  const success_url = `${process.env.NEXT_PUBLIC_SITE_URL}/billing/success`;
  const cancel_url = `${process.env.NEXT_PUBLIC_SITE_URL}/billing/cancel`;
  const checkout = await stripe.checkout.sessions.create({
    mode: "subscription",
    customer_email: user?.email || undefined,
    line_items: [{ price: STRIPE_PRICE_PRO, quantity: 1 }],
    success_url, cancel_url,
  });
  return NextResponse.json({ url: checkout.url });
}
