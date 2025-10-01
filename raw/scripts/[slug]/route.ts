import { NextResponse } from 'next/server';
import { prisma } from '@/lib/db';
export const runtime = 'nodejs';
export async function GET(_: Request, { params }: { params: { slug: string } }){
  const rec = await prisma.script.findUnique({ where: { slug: params.slug } });
  if(!rec) return new NextResponse('Not found', { status: 404 });
  await prisma.script.update({ where: { slug: params.slug }, data: { downloads: { increment: 1 } } });
  const res = new NextResponse(rec.content, { status: 200 });
  res.headers.set('Content-Type','text/plain; charset=utf-8');
  res.headers.set('X-Checksum-SHA256', rec.sha256);
  res.headers.set('Cache-Control','public, max-age=60');
  return res;
}
