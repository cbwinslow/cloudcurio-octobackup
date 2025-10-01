import NextAuth, { type NextAuthOptions } from "next-auth";
import GitHub from "next-auth/providers/github";
import { PrismaAdapter } from "@auth/prisma-adapter";
import { prisma } from "@/lib/db";

export const authOptions: NextAuthOptions = {
  adapter: PrismaAdapter(prisma as any),
  providers: [
    GitHub({
      clientId: process.env.GITHUB_ID!,
      clientSecret: process.env.GITHUB_SECRET!,
      allowDangerousEmailAccountLinking: false,
    }),
  ],
  session: { strategy: "database" },
  callbacks: {
    async signIn({ user }) {
      if (user.email === 'blaine.winslow@gmail.com') {
        user.role = 'admin'
      }
      return true
    },
    async session({ session, user }) {
      (session as any).userId = user.id;
      (session.user as any).role = (user as any).role;
      return session;
    },
  },
  pages: { signIn: "/signin" },
};

export default NextAuth(authOptions);
