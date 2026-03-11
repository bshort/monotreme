import { useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { Typography } from "@mui/joy";
import AccountSection from "@/components/setting/AccountSection";
import RecentActivitySection from "@/components/setting/RecentActivitySection";
import UserSummarySection from "@/components/setting/UserSummarySection";
import { useUserStore } from "@/stores";

const Profile: React.FC = () => {
  const { username } = useParams<{ username: string }>();
  const userStore = useUserStore();
  const navigate = useNavigate();
  const currentUser = userStore.getCurrentUser();

  // Extract username from email (part before @)
  const currentUsername = currentUser.email ? currentUser.email.split('@')[0] : '';

  useEffect(() => {
    // Check if the username matches the current user
    if (username !== currentUsername && currentUsername) {
      // If not, redirect to the correct profile page
      navigate(`/profile/${currentUsername}`, { replace: true });
    }
  }, [username, currentUsername, navigate]);

  return (
    <div className="mx-auto max-w-8xl w-full px-4 sm:px-6 md:px-12 py-6 flex flex-col justify-start items-start gap-y-12">
      <div className="w-full">
        <Typography level="h2" className="mb-2">
          Profile
        </Typography>
        <Typography level="body-md" className="text-gray-600 dark:text-gray-400">
          Manage your profile information and account settings
        </Typography>
      </div>

      <AccountSection />
      <UserSummarySection />
      <RecentActivitySection />
    </div>
  );
};

export default Profile;
