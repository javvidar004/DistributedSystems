// app/dashboard/page.tsx
// A concrete route that re-uses the dashboard implementation which lives
// inside the (main) grouping. The `(main)` folder is a route-group and
// isn't addressable directly via URL, so we expose `/dashboard` here and
// render the same UI.

import DashboardPage from '../(main)/page';
import MainLayout from '../(main)/layout';

export default function DashboardRoute() {
  // The (main) folder uses a route-group layout which isn't applied when
  // rendering the dashboard from a different path. Wrap the dashboard
  // component with the `MainLayout` so Navbar and sidebars are present.
  return (
    <MainLayout>
      <DashboardPage />
    </MainLayout>
  );
}
