import type { Metadata } from "next";
import ClientLayout from './client-layout';

export const metadata: Metadata = {
  title: "Frost Wallet",
  description: "Secure wallet setup",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>
        <ClientLayout>{children}</ClientLayout>
      </body>
    </html>
  );
}
