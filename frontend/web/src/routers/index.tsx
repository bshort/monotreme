import { createBrowserRouter } from "react-router-dom";
import App from "@/App";
import Root from "@/layouts/Root";
import AdminSignIn from "@/pages/AdminSignIn";
import AuthCallback from "@/pages/AuthCallback";
import BookmarkImport from "@/pages/BookmarkImport";
import CollectionDashboard from "@/pages/CollectionDashboard";
import CollectionSpace from "@/pages/CollectionSpace";
import Home from "@/pages/Home";
import NotFound from "@/pages/NotFound";
import ShortcutDashboard from "@/pages/ShortcutDashboard";
import ShortcutDetail from "@/pages/ShortcutDetail";
import ShortcutSpace from "@/pages/ShortcutSpace";
import SignIn from "@/pages/SignIn";
import SignUp from "@/pages/SignUp";
//import SubscriptionSetting from "@/pages/SubscriptionSetting";
import UserSetting from "@/pages/UserSetting";
import WorkspaceSetting from "@/pages/WorkspaceSetting";

const router = createBrowserRouter([
  {
    path: "/",
    element: <App />,
    children: [
      {
        path: "/auth",
        children: [
          {
            path: "",
            element: <SignIn />,
          },
          {
            path: "admin",
            element: <AdminSignIn />,
          },
          {
            path: "signup",
            element: <SignUp />,
          },
          {
            path: "callback",
            element: <AuthCallback />,
          },
        ],
      },
      {
        path: "",
        element: <Root />,
        children: [
          {
            path: "/",
            element: <Home />,
          },
          {
            path: "/shortcuts",
            element: <ShortcutDashboard />,
          },
          {
            path: "/collections",
            element: <CollectionDashboard />,
          },
          {
            path: "/shortcut/:shortcutId",
            element: <ShortcutDetail />,
          },
          {
            path: "/setting/general",
            element: <UserSetting />,
          },
          {
            path: "/setting/workspace",
            element: <WorkspaceSetting />,
          },
          {
            path: "/admin/import",
            element: <BookmarkImport />,
          },
        ],
      },
      {
        path: "s/*",
        element: <ShortcutSpace />,
      },
      {
        path: "c/*",
        element: <CollectionSpace />,
      },
      {
        path: "*",
        element: <NotFound />,
      },
    ],
  },
]);

export default router;
