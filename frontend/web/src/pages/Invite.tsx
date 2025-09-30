import { Button, Card, Typography } from "@mui/joy";
import { useEffect, useState } from "react";
import { toast } from "react-hot-toast";
import { useTranslation } from "react-i18next";
import { useSearchParams, Link } from "react-router-dom";
import Icon from "@/components/Icon";
import { authServiceClient } from "@/grpcweb";

const Invite: React.FC = () => {
  const { t } = useTranslation();
  const [searchParams] = useSearchParams();
  const [invitationCode, setInvitationCode] = useState<string>("");
  const [loading, setLoading] = useState(false);
  const [formData, setFormData] = useState({
    email: "",
    nickname: "",
    password: "",
    confirmPassword: "",
  });

  useEffect(() => {
    const code = searchParams.get("code");
    if (code) {
      setInvitationCode(code);
    }
  }, [searchParams]);

  const handleSignUp = async (e: React.FormEvent) => {
    e.preventDefault();

    if (formData.password !== formData.confirmPassword) {
      toast.error("Passwords do not match");
      return;
    }

    if (!formData.email || !formData.nickname || !formData.password) {
      toast.error("Please fill in all fields");
      return;
    }

    setLoading(true);
    try {
      await authServiceClient.signUp({
        email: formData.email,
        nickname: formData.nickname,
        password: formData.password,
        invitationCode: invitationCode,
      });
      window.location.href = "/";
    } catch (error: any) {
      console.error("Failed to sign up:", error);
      toast.error(error.details || "Failed to sign up");
    } finally {
      setLoading(false);
    }
  };

  const handleInputChange = (field: string, value: string) => {
    setFormData(prev => ({ ...prev, [field]: value }));
  };

  return (
    <div className="w-full h-screen flex flex-col justify-center items-center bg-gray-50 dark:bg-black">
      <div className="w-full max-w-md mx-auto px-4">
        <Card className="w-full p-6">
          <div className="flex flex-col items-center space-y-4">
            <div className="flex items-center space-x-2 mb-4">
              <Icon.UserPlus className="w-8 h-8" />
              <Typography level="h2">Join Monotreme</Typography>
            </div>

            {invitationCode && (
              <div className="w-full p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
                <div className="flex items-center space-x-2">
                  <Icon.Mail className="w-4 h-4 text-blue-600 dark:text-blue-400" />
                  <Typography level="body-sm" className="text-blue-800 dark:text-blue-200">
                    You've been invited to join Monotreme! Complete your registration below.
                  </Typography>
                </div>
              </div>
            )}

            <form onSubmit={handleSignUp} className="w-full space-y-4">
              <div>
                <Typography level="body-sm" className="mb-1 font-medium">
                  Email
                </Typography>
                <input
                  type="email"
                  value={formData.email}
                  onChange={(e) => handleInputChange("email", e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="Enter your email"
                  required
                />
              </div>

              <div>
                <Typography level="body-sm" className="mb-1 font-medium">
                  Nickname
                </Typography>
                <input
                  type="text"
                  value={formData.nickname}
                  onChange={(e) => handleInputChange("nickname", e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="Enter your nickname"
                  required
                />
              </div>

              <div>
                <Typography level="body-sm" className="mb-1 font-medium">
                  Password
                </Typography>
                <input
                  type="password"
                  value={formData.password}
                  onChange={(e) => handleInputChange("password", e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="Enter your password"
                  required
                />
              </div>

              <div>
                <Typography level="body-sm" className="mb-1 font-medium">
                  Confirm Password
                </Typography>
                <input
                  type="password"
                  value={formData.confirmPassword}
                  onChange={(e) => handleInputChange("confirmPassword", e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="Confirm your password"
                  required
                />
              </div>

              <Button
                type="submit"
                disabled={loading}
                className="w-full"
                size="lg"
              >
                {loading ? (
                  <div className="flex items-center space-x-2">
                    <Icon.Loader className="w-4 h-4 animate-spin" />
                    <span>Creating account...</span>
                  </div>
                ) : (
                  <div className="flex items-center space-x-2">
                    <Icon.UserPlus className="w-4 h-4" />
                    <span>Create Account</span>
                  </div>
                )}
              </Button>
            </form>

            <div className="w-full pt-4 border-t border-gray-200 dark:border-zinc-700">
              <Typography level="body-sm" className="text-center text-gray-600 dark:text-gray-400">
                Already have an account?{" "}
                <Link
                  to="/auth"
                  className="text-blue-600 dark:text-blue-400 hover:underline"
                >
                  Sign in
                </Link>
              </Typography>
            </div>
          </div>
        </Card>
      </div>
    </div>
  );
};

export default Invite;