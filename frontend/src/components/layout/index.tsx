// src/components/Layout.tsx
import { Outlet } from 'react-router';

export function Layout() {
  return (
    <div>
      {/* Shared UI, e.g. sidebar/header */}
      <header>My App Header</header>
      <main>
        <Outlet />
      </main>
    </div>
  );
}