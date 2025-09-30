import { Button, Card, Typography } from "@mui/joy";
import { useEffect, useState } from "react";
import { QRCodeSVG } from "qrcode.react";
import { toast } from "react-hot-toast";
import { useTranslation } from "react-i18next";
import Icon from "@/components/Icon";
import { userServiceClient } from "@/grpcweb";
import { useUserStore } from "@/stores";

const QRCode: React.FC = () => {
  const { t } = useTranslation();
  const userStore = useUserStore();
  const currentUser = userStore.getCurrentUser();
  const [invitationCode, setInvitationCode] = useState<string>("");
  const [invitationUrl, setInvitationUrl] = useState<string>("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchInvitationCode = async () => {
      if (!currentUser || currentUser.id === -1) {
        setLoading(false);
        return;
      }

      try {
        const response = await userServiceClient.getInvitationCode({});
        setInvitationCode(response.invitationCode);
        setInvitationUrl(window.location.origin + response.invitationUrl);
      } catch (error: any) {
        console.error("Failed to get invitation code:", error);
        toast.error("Failed to load QR code");
      } finally {
        setLoading(false);
      }
    };

    fetchInvitationCode();
  }, [currentUser]);

  const handleCopyUrl = () => {
    if (invitationUrl) {
      navigator.clipboard.writeText(invitationUrl);
      toast.success("Invitation URL copied to clipboard");
    }
  };

  const handleCopyCode = () => {
    if (invitationCode) {
      navigator.clipboard.writeText(invitationCode);
      toast.success("Invitation code copied to clipboard");
    }
  };

  if (!currentUser || currentUser.id === -1) {
    return (
      <div className="w-full max-w-8xl mx-auto px-4 sm:px-6 md:px-12 py-6">
        <div className="w-full flex flex-col justify-start items-start">
          <div className="w-full flex flex-row justify-start items-center mb-6">
            <Icon.QrCode className="w-6 h-auto mr-2" />
            <Typography level="h2">QR Code for Invitations</Typography>
          </div>
          <Card className="w-full max-w-md mx-auto">
            <div className="flex justify-center items-center py-8">
              <Typography level="body-md" className="text-center text-gray-600 dark:text-gray-400">
                Please sign in to generate invitation QR codes.
              </Typography>
            </div>
          </Card>
        </div>
      </div>
    );
  }

  if (loading) {
    return (
      <div className="w-full max-w-8xl mx-auto px-4 sm:px-6 md:px-12 py-6">
        <div className="w-full flex flex-col justify-start items-start">
          <div className="w-full flex flex-row justify-start items-center mb-6">
            <Icon.QrCode className="w-6 h-auto mr-2" />
            <Typography level="h2">QR Code for Invitations</Typography>
          </div>
          <Card className="w-full max-w-md mx-auto">
            <div className="flex justify-center items-center py-8">
              <Icon.Loader className="w-6 h-6 animate-spin" />
              <Typography level="body-md" className="ml-2">
                Loading QR code...
              </Typography>
            </div>
          </Card>
        </div>
      </div>
    );
  }

  return (
    <div className="w-full max-w-8xl mx-auto px-4 sm:px-6 md:px-12 py-6">
      <div className="w-full flex flex-col justify-start items-start">
        <div className="w-full flex flex-row justify-start items-center mb-6">
          <Icon.QrCode className="w-6 h-auto mr-2" />
          <Typography level="h2">QR Code for Invitations</Typography>
        </div>

        <Card className="w-full max-w-md mx-auto p-6">
          <div className="flex flex-col items-center space-y-4">
            <Typography level="body-md" className="text-center text-gray-600 dark:text-gray-400">
              Share this QR code to invite new users to sign up. They'll be marked as invited by{" "}
              <strong>{currentUser.nickname}</strong>.
            </Typography>

            <div className="bg-white p-4 rounded-lg">
              <QRCodeSVG
                value={invitationUrl}
                size={200}
                level="M"
                includeMargin={true}
              />
            </div>

            <div className="w-full space-y-2">
              <div className="flex flex-col space-y-1">
                <Typography level="body-sm" className="font-medium">
                  Invitation URL:
                </Typography>
                <div className="flex items-center space-x-2">
                  <Typography
                    level="body-sm"
                    className="flex-1 px-2 py-1 bg-gray-100 dark:bg-zinc-800 rounded text-sm font-mono break-all"
                  >
                    {invitationUrl}
                  </Typography>
                  <Button size="sm" variant="outlined" onClick={handleCopyUrl}>
                    <Icon.Copy className="w-4 h-4" />
                  </Button>
                </div>
              </div>

              <div className="flex flex-col space-y-1">
                <Typography level="body-sm" className="font-medium">
                  Invitation Code:
                </Typography>
                <div className="flex items-center space-x-2">
                  <Typography
                    level="body-sm"
                    className="flex-1 px-2 py-1 bg-gray-100 dark:bg-zinc-800 rounded text-sm font-mono"
                  >
                    {invitationCode}
                  </Typography>
                  <Button size="sm" variant="outlined" onClick={handleCopyCode}>
                    <Icon.Copy className="w-4 h-4" />
                  </Button>
                </div>
              </div>
            </div>

            <Typography level="body-xs" className="text-center text-gray-500 dark:text-gray-400">
              New users can visit the invitation URL or use the code during signup to be linked to your account.
            </Typography>
          </div>
        </Card>
      </div>
    </div>
  );
};

export default QRCode;