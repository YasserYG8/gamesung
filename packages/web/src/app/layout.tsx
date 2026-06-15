import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "gamesung",
  description: "An open-source platform for creating, building, and sharing games",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
