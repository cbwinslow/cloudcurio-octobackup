import { prisma } from '@/lib/db';
import { notFound } from 'next/navigation';
export default async function ReviewPage({ params }:{ params:{ id:string } }){
  const art = await prisma.reviewArtifact.findUnique({ where: { jobId: params.id }, include: { job: true } });
  if(!art) return notFound();
  return (
    <main className="max-w-5xl mx-auto p-6">
      <h1 className="text-2xl font-bold mb-4">Review Report</h1>
      <article className="prose" dangerouslySetInnerHTML={{ __html: art.content }} />
    </main>
  );
}
