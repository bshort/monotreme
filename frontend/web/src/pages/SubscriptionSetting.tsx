import { Alert, Button, Divider, Link, Textarea } from "@mui/joy";
import { useState } from "react";
import toast from "react-hot-toast";
import { showCommonDialog } from "@/components/Alert";
import Icon from "@/components/Icon";
import SubscriptionFAQ from "@/components/SubscriptionFAQ";
import { subscriptionServiceClient } from "@/grpcweb";
import { useUserStore, useWorkspaceStore } from "@/stores";
import { stringifyPlanType } from "@/stores/subscription";
import { PlanType } from "@/types/proto/api/v1/subscription_service";
import { Role } from "@/types/proto/api/v1/user_service";

const SubscriptionSetting: React.FC = () => {
  const workspaceStore = useWorkspaceStore();
  const currentUser = useUserStore().getCurrentUser();
  const [licenseKey, setLicenseKey] = useState<string>("");
  const isAdmin = currentUser.role === Role.ADMIN;
  const subscription = workspaceStore.getSubscription();

  const handleDeleteLicenseKey = async () => {
    if (!isAdmin) {
      toast.error("Only admin can upload license key");
      return;
    }

    showCommonDialog({
      title: "Reset licence key",
      content: `Are you sure to reset the license key? You cannot undo this action.`,
      style: "warning",
      onConfirm: async () => {
        try {
          await subscriptionServiceClient.deleteSubscription({});
          toast.success("License key has been reset");
        } catch (error: any) {
          toast.error(error.details);
        }
        await workspaceStore.fetchWorkspaceProfile();
      },
    });
  };
/*
  const handleUpdateLicenseKey = async () => {
    if (!isAdmin) {
      toast.error("Only admin can upload license key");
      return;
    }

    try {
      const subscription = await subscriptionServiceClient.updateSubscription({
        licenseKey,
      });
      toast.success(`Welcome to Monotreme ${stringifyPlanType(subscription.plan)}🎉`);
    } catch (error: any) {
      toast.error(error.details);
    }
    setLicenseKey("");
    await workspaceStore.fetchWorkspaceProfile();
    */
  };

  return (
    <div className="mx-auto max-w-8xl w-full px-4 sm:px-6 md:px-12 pt-8 pb-24 flex flex-col justify-start items-start gap-y-12">
      <div className="w-full">
        <p className="text-2xl shrink-0 font-semibold text-gray-900 dark:text-gray-500">Subscription</p>
        <div className="mt-2">
          <span className="text-gray-500 mr-2">Current plan:</span>
          <span className="text-2xl mr-4 dark:text-gray-400">{stringifyPlanType(subscription.plan)}</span>
        </div>
        <Textarea
          className="w-full mt-2"
          minRows={2}
          maxRows={2}
          placeholder="Enter your license key here - write only"
          value={licenseKey}
          onChange={(event) => setLicenseKey(event.target.value)}
        />
        <div className="w-full flex justify-between items-center mt-4">
          <div>
            {subscription.plan === PlanType.FREE && (
              <Link href="https://monotrememarks.lemonsqueezy.com/checkout/buy/947e9a56-c93a-4294-8d71-2ea4b0f3ec51" target="_blank">
                Buy a license key
                <Icon.ExternalLink className="w-4 h-auto ml-1" />
              </Link>
            )}
          </div>
          <div className="flex justify-end items-center gap-2">
            {subscription.plan !== PlanType.FREE && (
              <Button color="neutral" variant="plain" onClick={handleDeleteLicenseKey}>
                Reset
              </Button>
            )}
            <Button disabled={licenseKey === ""} onClick={handleUpdateLicenseKey}>
              Upload license
            </Button>
          </div>
        </div>
      </div>
      <Divider />
      <section className="w-full pb-8 dark:bg-zinc-900 flex items-center justify-center">
        <div className="w-full px-6">
          <div className="max-w-4xl mx-auto mb-12">
            <Alert className="!inline-block mb-12">
              Monotreme is an open source, self-hosted platform for sharing and managing your most frequently used links. Easily create
              customizable, human-readable shortcuts to streamline your link management. Our source code is available and accessible on{" "}
              <Link href="https://github.com/bshort/monotreme" target="_blank">
                GitHub
              </Link>{" "}
              so anyone can get it, inspect it and review it.
            </Alert>
          </div>
          <div className="w-full grid grid-cols-1 gap-6 lg:gap-12 mt-8 md:grid-cols-3 md:max-w-4xl mx-auto">
            <div className="flex flex-col p-6 bg-white dark:bg-zinc-800 shadow-lg rounded-lg justify-between border border-gray-300 dark:border-zinc-700">
              <div>
                <h3 className="text-2xl font-bold text-center dark:text-gray-300">Free</h3>
                <div className="mt-3 text-center text-zinc-600 dark:text-zinc-400">
                  <span className="text-4xl font-bold">$0</span>/ month
                </div>
                <ul className="mt-4 space-y-3">
                  <li className="flex items-center dark:text-gray-300">
                    <Icon.CheckCircle2 className="w-5 h-auto text-green-600 mr-1 shrink-0" />
                    Full API access
                  </li>
                  <li className="flex items-center dark:text-gray-300">
                    <Icon.CheckCircle2 className="w-5 h-auto text-green-600 mr-1 shrink-0" />
                    Browser extension
                  </li>
                  <li className="flex items-center dark:text-gray-300">
                    <Icon.CheckCircle2 className="w-5 h-auto text-green-600 mr-1 shrink-0" />
                    Basic analytics
                  </li>
                  <li className="flex items-center dark:text-gray-300">
                    <Icon.AlertCircle className="w-5 h-auto text-gray-400 mr-1 shrink-0" />
                    Up to 100 shortcuts
                  </li>
                  <li className="flex items-center dark:text-gray-300">
                    <Icon.AlertCircle className="w-5 h-auto text-gray-400 mr-1 shrink-0" />
                    Up to 5 collections
                  </li>
                  <li className="flex items-center dark:text-gray-300">
                    <Icon.AlertCircle className="w-5 h-auto text-gray-400 mr-1 shrink-0" />
                    Up to 5 members
                  </li>
                </ul>
              </div>
            </div>
            <div className="relative flex flex-col p-6 bg-white dark:bg-zinc-800 shadow-lg rounded-lg dark:bg-zinc-850 justify-between border border-purple-500">
              <div className="px-3 py-1 text-sm text-white bg-gradient-to-r from-pink-500 to-purple-500 rounded-full inline-block absolute top-0 left-1/2 transform -translate-x-1/2 -translate-y-1/2">
                Popular
              </div>
              <div>
                <h3 className="text-2xl font-bold text-center dark:text-gray-300">Pro</h3>
                <div className="mt-3 text-center text-zinc-600 dark:text-zinc-400">
                  <span className="text-4xl font-bold">$4</span>/ month
                </div>
                <p className="mt-3 font-medium dark:text-gray-300">Everything in Free, and</p>
                <ul className="mt-4 space-y-3">
                  <li className="flex items-center dark:text-gray-300">
                    <Icon.CheckCircle2 className="w-5 h-auto text-green-600 mr-1 shrink-0" />
                    Unlimited members
                  </li>
                  <li className="flex items-center dark:text-gray-300">
                    <Icon.CheckCircle2 className="w-5 h-auto text-green-600 mr-1 shrink-0" />
                    Unlimited shortcuts
                  </li>
                  <li className="flex items-center dark:text-gray-300">
                    <Icon.CheckCircle2 className="w-5 h-auto text-green-600 mr-1 shrink-0" />
                    Unlimited collections
                  </li>
                  <li className="flex items-center dark:text-gray-300">
                    <Icon.CheckCircle2 className="w-5 h-auto text-green-600 mr-1 shrink-0" />
                    High-priority in roadmap
                  </li>
                  <li className="flex items-center dark:text-gray-300">
                    <Icon.CheckCircle2 className="w-5 h-auto text-green-600 mr-1 shrink-0" />
                    Email support
                  </li>
                </ul>
              </div>
              <div className="mt-6">
                <Link
                  className="w-full"
                  underline="none"
                  href="https://monotrememarks.lemonsqueezy.com/checkout/buy/947e9a56-c93a-4294-8d71-2ea4b0f3ec51"
                  target="_blank"
                >
                  <Button className="w-full bg-gradient-to-r from-pink-500 to-purple-500 shadow hover:opacity-80">Get Pro License</Button>
                </Link>
              </div>
            </div>
            <div className="flex flex-col p-6 bg-white dark:bg-zinc-800 shadow-lg rounded-lg dark:bg-zinc-850 justify-between border border-gray-300 dark:border-zinc-700">
              <div>
                <span className="block text-2xl text-center dark:text-gray-200 opacity-80">Team</span>
                <div className="mt-3 text-center text-zinc-600 dark:text-zinc-400">
                  <span className="mr-2">start with</span>
                  <span>
                    <span className="text-4xl font-bold">$10</span>/ month
                  </span>
                </div>
                <p className="mt-3 font-medium dark:text-gray-300">Everything in Pro, and</p>
                <ul className="mt-4 space-y-3">
                  <li className="flex items-center dark:text-gray-300">
                    <Icon.BadgeCheck className="w-5 h-auto text-green-600 mr-1 shrink-0" />
                    Custom branding
                  </li>
                  <li className="flex items-center dark:text-gray-300">
                    <Icon.BarChart3 className="w-5 h-auto text-green-600 mr-1 shrink-0" />
                    Advanced analytics
                  </li>
                  <li className="flex items-center dark:text-gray-300">
                    <Icon.Shield className="w-5 h-auto text-green-600 mr-1 shrink-0" />
                    Single Sign-On(SSO)
                  </li>
                  <li className="flex items-center dark:text-gray-300">
                    <Icon.Building2 className="w-5 h-auto text-green-600 mr-1 shrink-0" />
                    Priority support
                  </li>
                  <li className="flex items-center dark:text-gray-300">
                    <Icon.Sparkles className="w-5 h-auto text-green-600 mr-1 shrink-0" />
                    More coming soon
                  </li>
                </ul>
              </div>
              <div className="mt-6">
                <Link
                  className="w-full"
                  underline="none"
                  href="mailto:test@gmail.com?subject=Inquiry about Monotreme Team Plan"
                  target="_blank"
                >
                  <Button className="w-full">Contact us</Button>
                </Link>
              </div>
            </div>
          </div>
        </div>
      </section>
      <SubscriptionFAQ />
    </div>
  );
};

export default SubscriptionSetting;
