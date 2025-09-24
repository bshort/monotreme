import { Button, Input } from "@mui/joy";
import React, { FormEvent, useEffect, useState } from "react";
import { toast } from "react-hot-toast";
import { useTranslation } from "react-i18next";
import { Link } from "react-router-dom";
import Logo from "@/components/Logo";
import { authServiceClient } from "@/grpcweb";
import useLoading from "@/hooks/useLoading";
import useNavigateTo from "@/hooks/useNavigateTo";
import { useUserStore, useWorkspaceStore } from "@/stores";

const Landing: React.FC = () => {
  const { t } = useTranslation();
  const navigateTo = useNavigateTo();
  const workspaceStore = useWorkspaceStore();
  const userStore = useUserStore();
  const [email, setEmail] = useState("");
  const [nickname, setNickname] = useState("");
  const [password, setPassword] = useState("");
  const actionBtnLoadingState = useLoading(false);
  const allowConfirm = email.length > 0 && nickname.length > 0 && password.length > 0;

  useEffect(() => {
    if (workspaceStore.setting.disallowUserRegistration) {
      return navigateTo("/auth", {
        replace: true,
      });
    }
  }, []);

  const handleEmailInputChanged = (e: React.ChangeEvent<HTMLInputElement>) => {
    const text = e.target.value as string;
    setEmail(text);
  };

  const handleNicknameInputChanged = (e: React.ChangeEvent<HTMLInputElement>) => {
    const text = e.target.value as string;
    setNickname(text);
  };

  const handlePasswordInputChanged = (e: React.ChangeEvent<HTMLInputElement>) => {
    const text = e.target.value as string;
    setPassword(text);
  };

  const handleSignupBtnClick = async (e: FormEvent) => {
    e.preventDefault();
    if (actionBtnLoadingState.isLoading) {
      return;
    }

    try {
      actionBtnLoadingState.setLoading();
      const user = await authServiceClient.signUp({
        email,
        nickname,
        password,
      });
      if (user) {
        userStore.setCurrentUserId(user.id);
        await userStore.fetchCurrentUser();
        navigateTo("/");
      } else {
        toast.error("Signup failed");
      }
    } catch (error: any) {
      console.error(error);
      toast.error(error.details);
    }
    actionBtnLoadingState.setFinish();
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-indigo-50 dark:from-zinc-900 dark:via-zinc-800 dark:to-zinc-900">
      {/* Header */}
      <header className="w-full px-6 py-4">
        <div className="flex items-center justify-between max-w-6xl mx-auto">
          <div className="flex items-center space-x-3">
            <Logo className="w-8 h-8" />
            <span className="text-2xl font-bold text-gray-900 dark:text-white">Monotreme</span>
          </div>
          <Link
            to="/auth"
            className="text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300 font-medium transition-colors"
          >
            Sign In
          </Link>
        </div>
      </header>

      {/* Hero Section */}
      <section className="w-full px-6 py-16">
        <div className="max-w-6xl mx-auto text-center">
          <h1 className="text-5xl md:text-6xl font-bold text-gray-900 dark:text-white mb-6">
            Your Personal
            <span className="text-blue-600 dark:text-blue-400"> Link Management</span>
            <br />
            Made Simple
          </h1>
          <p className="text-xl text-gray-600 dark:text-gray-300 mb-8 max-w-3xl mx-auto leading-relaxed">
            Monotreme helps you organize, share, and track your links with ease. Create custom shortcuts, organize them into collections, and access them from anywhere.
          </p>
        </div>
      </section>

      {/* Features Section */}
      <section className="w-full px-6 py-16 bg-white/50 dark:bg-zinc-800/50">
        <div className="max-w-6xl mx-auto">
          <div className="grid md:grid-cols-3 gap-8">
            <div className="text-center p-6">
              <div className="w-16 h-16 bg-blue-100 dark:bg-blue-900/30 rounded-xl flex items-center justify-center mx-auto mb-4">
                <svg className="w-8 h-8 text-blue-600 dark:text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
                </svg>
              </div>
              <h3 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">Smart Shortcuts</h3>
              <p className="text-gray-600 dark:text-gray-300">Create memorable short links that are easy to share and remember. Perfect for frequently accessed resources.</p>
            </div>

            <div className="text-center p-6">
              <div className="w-16 h-16 bg-green-100 dark:bg-green-900/30 rounded-xl flex items-center justify-center mx-auto mb-4">
                <svg className="w-8 h-8 text-green-600 dark:text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                </svg>
              </div>
              <h3 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">Organized Collections</h3>
              <p className="text-gray-600 dark:text-gray-300">Group your links into collections and use tags to keep everything organized and easily searchable.</p>
            </div>

            <div className="text-center p-6">
              <div className="w-16 h-16 bg-purple-100 dark:bg-purple-900/30 rounded-xl flex items-center justify-center mx-auto mb-4">
                <svg className="w-8 h-8 text-purple-600 dark:text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                </svg>
              </div>
              <h3 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">Track & Analyze</h3>
              <p className="text-gray-600 dark:text-gray-300">Monitor click counts and usage patterns to understand how your links are being used.</p>
            </div>
          </div>
        </div>
      </section>

      {/* Benefits Section */}
      <section className="w-full px-6 py-16">
        <div className="max-w-6xl mx-auto">
          <h2 className="text-3xl md:text-4xl font-bold text-center text-gray-900 dark:text-white mb-12">
            Why Choose Monotreme?
          </h2>
          <div className="grid md:grid-cols-2 gap-12 items-center">
            <div>
              <ul className="space-y-4">
                <li className="flex items-start space-x-3">
                  <div className="w-6 h-6 bg-blue-100 dark:bg-blue-900/30 rounded-full flex items-center justify-center flex-shrink-0 mt-0.5">
                    <svg className="w-4 h-4 text-blue-600 dark:text-blue-400" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                    </svg>
                  </div>
                  <div>
                    <h4 className="font-semibold text-gray-900 dark:text-white">Self-Hosted & Private</h4>
                    <p className="text-gray-600 dark:text-gray-300">Keep your data under your control with self-hosting capabilities.</p>
                  </div>
                </li>
                <li className="flex items-start space-x-3">
                  <div className="w-6 h-6 bg-blue-100 dark:bg-blue-900/30 rounded-full flex items-center justify-center flex-shrink-0 mt-0.5">
                    <svg className="w-4 h-4 text-blue-600 dark:text-blue-400" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                    </svg>
                  </div>
                  <div>
                    <h4 className="font-semibold text-gray-900 dark:text-white">Team Collaboration</h4>
                    <p className="text-gray-600 dark:text-gray-300">Share collections with team members and collaborate on link organization.</p>
                  </div>
                </li>
                <li className="flex items-start space-x-3">
                  <div className="w-6 h-6 bg-blue-100 dark:bg-blue-900/30 rounded-full flex items-center justify-center flex-shrink-0 mt-0.5">
                    <svg className="w-4 h-4 text-blue-600 dark:text-blue-400" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                    </svg>
                  </div>
                  <div>
                    <h4 className="font-semibold text-gray-900 dark:text-white">Browser Extensions</h4>
                    <p className="text-gray-600 dark:text-gray-300">Quickly save bookmarks with our Chrome extension for seamless workflow.</p>
                  </div>
                </li>
                <li className="flex items-start space-x-3">
                  <div className="w-6 h-6 bg-blue-100 dark:bg-blue-900/30 rounded-full flex items-center justify-center flex-shrink-0 mt-0.5">
                    <svg className="w-4 h-4 text-blue-600 dark:text-blue-400" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                    </svg>
                  </div>
                  <div>
                    <h4 className="font-semibold text-gray-900 dark:text-white">Multiple View Options</h4>
                    <p className="text-gray-600 dark:text-gray-300">Switch between card, list, and compact views to match your preference.</p>
                  </div>
                </li>
              </ul>
            </div>
            
            {/* Signup Form */}
            <div className="bg-white dark:bg-zinc-800 rounded-2xl shadow-xl p-8">
              <h3 className="text-2xl font-bold text-gray-900 dark:text-white mb-6 text-center">
                Get Started Today
              </h3>
              <form className="space-y-4" onSubmit={handleSignupBtnClick}>
                <div className={`space-y-4 ${actionBtnLoadingState.isLoading ? "opacity-80" : ""}`}>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Email Address
                    </label>
                    <Input
                      type="email"
                      value={email}
                      placeholder="your@email.com"
                      onChange={handleEmailInputChanged}
                      className="w-full"
                      disabled={actionBtnLoadingState.isLoading}
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Display Name
                    </label>
                    <Input
                      type="text"
                      value={nickname}
                      placeholder="Your Name"
                      onChange={handleNicknameInputChanged}
                      className="w-full"
                      disabled={actionBtnLoadingState.isLoading}
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                      Password
                    </label>
                    <Input
                      type="password"
                      value={password}
                      placeholder="••••••••"
                      onChange={handlePasswordInputChanged}
                      className="w-full"
                      disabled={actionBtnLoadingState.isLoading}
                    />
                  </div>
                </div>
                
                <Button
                  type="submit"
                  color="primary"
                  size="lg"
                  loading={actionBtnLoadingState.isLoading}
                  disabled={actionBtnLoadingState.isLoading || !allowConfirm}
                  className="w-full"
                  onClick={handleSignupBtnClick}
                >
                  Create Your Account
                </Button>
              </form>
              
              {workspaceStore.profile.owner && (
                <p className="text-center text-sm text-gray-600 dark:text-gray-400 mt-4">
                  Already have an account?{" "}
                  <Link to="/auth" className="text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300 font-medium">
                    Sign In
                  </Link>
                </p>
              )}
            </div>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="w-full px-6 py-8 bg-gray-50 dark:bg-zinc-800 border-t border-gray-200 dark:border-zinc-700">
        <div className="max-w-6xl mx-auto text-center">
          <div className="flex items-center justify-center space-x-3 mb-4">
            <Logo className="w-6 h-6" />
            <span className="text-lg font-semibold text-gray-900 dark:text-white">Monotreme</span>
          </div>
          <p className="text-gray-600 dark:text-gray-400 text-sm">
            Your personal link management solution. Simple, powerful, and self-hosted.
          </p>
        </div>
      </footer>
    </div>
  );
};

export default Landing;