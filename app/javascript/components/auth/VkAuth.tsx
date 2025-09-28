import { router } from '@inertiajs/react';
import { Center } from '@mantine/core';
import * as VKID from '@vkid/sdk';
import i18next from 'i18next';
import React, { useEffect, useRef } from 'react';
import { generateCodeChallenge, generateCodeVerifier } from '@/lib/utils.ts';
import * as Routes from '@/routes.js';

interface VkIdLoginPayload {
  code: string;
  device_id: string;
  state: string;
}

interface VkIdError {
  code: number;
  text: string;
}

const langMap: Record<string, VKID.Languages> = {
  ru: VKID.Languages.RUS,
  en: VKID.Languages.ENG,
};

const handleVkAuth = (payload: VkIdLoginPayload) => {
  const verifier = localStorage.getItem('vk_code_verifier');
  const params = {
    ...payload,
    code_verifier: verifier,
  };
  router.get(Routes.callback_auth_path('vk'), params);
};

export default function VkAuth() {
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    (async () => {
      if (
        !ref.current ||
        !import.meta.env.VITE_VK_APP_ID ||
        !import.meta.env.VITE_VK_APP_REDIRECT_URL
      ) {
        return;
      }

      const codeVerifier = generateCodeVerifier();
      localStorage.setItem('vk_code_verifier', codeVerifier);
      const codeChallenge = await generateCodeChallenge(codeVerifier);
      const state = crypto.randomUUID();

      VKID.Config.init({
        app: import.meta.env.VITE_VK_APP_ID,
        redirectUrl: import.meta.env.VITE_VK_APP_REDIRECT_URL,
        codeChallenge: codeChallenge,
        state: state,
        responseMode: VKID.ConfigResponseMode.Callback,
        scope: 'email',
      });

      const oneTap = new VKID.OneTap();
      ref.current &&
        oneTap.render({
          container: ref.current,
          lang: langMap[i18next.language] || VKID.Languages.RUS,
        });
      oneTap.on(VKID.OneTapInternalEvents.LOGIN_SUCCESS, handleVkAuth);
      oneTap.on(VKID.WidgetEvents.ERROR, (error: VkIdError) => {
        throw new Error(error.text);
      });
    })();
  }, []);

  return <Center mt={10} ref={ref}></Center>;
}
