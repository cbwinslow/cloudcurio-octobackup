import { prisma } from "@/lib/db";
import { PLAN_LIMITS, getUserPlan } from "@/lib/entitlements";
export async function recordUsage(userId: string, kind: string, tokensIn=0, tokensOut=0, meta: any=null){
  await prisma.usageEvent.create({ data: { userId, kind, tokensIn, tokensOut, meta } });
}
export async function checkQuota(userId: string){
  const plan = await getUserPlan(userId);
  const today = new Date(); today.setHours(0,0,0,0);
  const count = await prisma.usageEvent.count({ where: { userId, kind: 'chat_request', createdAt: { gte: today } } });
  const limits = PLAN_LIMITS[plan];
  return { ok: count < limits.dailyRequests, used: count, limit: limits.dailyRequests, plan };
}
