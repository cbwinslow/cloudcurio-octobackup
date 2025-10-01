'use client';
import { signIn } from "next-auth/react";
export default function SignIn(){
  return (
    <main className="max-w-sm mx-auto py-24">
      <h1 className="text-2xl font-bold mb-4">Sign in</h1>
      <button className="w-full rounded bg-black text-white py-2" onClick={()=>signIn("github")}>Sign in with GitHub</button>
    </main>
  );
}
