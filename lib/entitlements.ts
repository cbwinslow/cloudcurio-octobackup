import { prisma } from "@/lib/db";
export type Plan = "free" | "pro" | "enterprise";
export const PLAN_LIMITS: Record<Plan, { dailyRequests: number; rpm: number; maxTokens: number; }>= {
  free: { dailyRequests: 50,  rpm: 10, maxTokens: 2000 },
  pro:  { dailyRequests: 500, rpm: 60, maxTokens: 8000 },
  enterprise: { dailyRequests: 5000, rpm: 120, maxTokens: 32000 },
};
export async function getUserPlan(userId: string): Promise<Plan> {
  const sub = await prisma.subscription.findFirst({ where: { userId, status: { in: ["active","trialing","past_due"] } } });
  if(!sub) return "free";
  if(sub.plan === 'enterprise') return 'enterprise';
  return 'pro';
}
