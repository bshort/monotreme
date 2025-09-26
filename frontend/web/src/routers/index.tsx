import { createBrowserRouter } from "react-router-dom";
import App from "@/App";
import Root from "@/layouts/Root";
import AdminSignIn from "@/pages/AdminSignIn";
import AuthCallback from "@/pages/AuthCallback";
import BookmarkImport from "@/pages/BookmarkImport";
import CollectionDashboard from "@/pages/CollectionDashboard";
import CollectionSpace from "@/pages/CollectionSpace";
import Home from "@/pages/Home";
import Landing from "@/pages/Landing";
import NotFound from "@/pages/NotFound";
import QuickSave from "@/pages/QuickSave";
import ShortcutDashboard from "@/pages/ShortcutDashboard";
import ShortcutDetail from "@/pages/ShortcutDetail";
import ShortcutSpace from "@/pages/ShortcutSpace";
import SignIn from "@/pages/SignIn";
import SignUp from "@/pages/SignUp";
import Stats from "@/pages/Stats";
//import SubscriptionSetting from "@/pages/SubscriptionSetting";
import TagsDashboard from "@/pages/TagsDashboard";
import UserSetting from "@/pages/UserSetting";
import WorkspaceSetting from "@/pages/WorkspaceSetting";

const router = createBrowserRouter([
  {
    path: "/",
    element: <App />,
    children: [
      {
        path: "/landing",
        element: <Landing />,
      },
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
            path: "/tags",
            element: <TagsDashboard />,
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
            path: "/stats",
            element: <Stats />,
          },
          {
            path: "/admin/import",
            element: <BookmarkImport />,
          },
        ],
      },
      {
        path: "/quick-save",
        element: <QuickSave />,
      },
      {
        path: "s/*",
        element: <ShortcutSpace />,
      },
      {
        path: ":prefix/*",
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
