import { prisma } from '@/lib/db';
export default async function Dotfiles(){
  const items = await prisma.script.findMany({ where: { channel: 'dotfile' }, orderBy: { updatedAt:'desc' } });
  return (
    <main className="max-w-4xl mx-auto p-6">
      <h1 className="text-3xl font-bold mb-4">Dotfiles</h1>
      <ul className="space-y-3">
        {items.map(s => (
          <li key={s.id} className="border rounded-xl p-3">
            <div className="font-mono">/{s.slug}</div>
            <pre className="bg-gray-50 p-3 rounded-xl overflow-auto text-xs">{`curl -fsSL https://cloudcurio.cc/raw/scripts/${s.slug} -o ~/.${s.slug}`}</pre>
          </li>
        ))}
      </ul>
    </main>
  );
}
