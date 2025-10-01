import { PrismaClient } from '@prisma/client'
import bcrypt from 'bcryptjs'  // npm install bcryptjs @types/bcryptjs if needed

const prisma = new PrismaClient()

async function seedAdmin() {
  const email = 'blaine.winslow@gmail.com'
  const name = 'cbwinslow'
  const hashedPassword = await bcrypt.hash('your_initial_password_here', 12)  // Change this!

  const user = await prisma.user.upsert({
    where: { email },
    update: {},
    create: {
      email,
      name,
      role: 'admin',  // Custom role if extended; default is member
      // For NextAuth, password not used (OAuth), but for direct auth if added
    },
  })

  console.log(`Admin user created/verified: ${user.name} (${user.email}) with role: ${user.role}`)

  await prisma.$disconnect()
}

seedAdmin()
  .catch((e) => {
    console.error(e)
    process.exit(1)
  })
  .finally(async () => {
    await prisma.$disconnect()
  })
