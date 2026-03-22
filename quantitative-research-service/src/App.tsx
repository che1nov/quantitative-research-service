import { useState, useEffect, ReactNode } from 'react';
import bridge, { UserInfo } from '@vkontakte/vk-bridge';
import { View, SplitLayout, SplitCol, ScreenSpinner } from '@vkontakte/vkui';
import { useActiveVkuiLocation } from '@vkontakte/vk-mini-apps-router';

import { Persik, Home, Surveys } from './panels';
import { DEFAULT_VIEW_PANELS } from './routes';

export const App = () => {
  const { panel: activePanel = DEFAULT_VIEW_PANELS.HOME } = useActiveVkuiLocation();
  const [fetchedUser, setUser] = useState<UserInfo | undefined>();
  const [authStatus, setAuthStatus] = useState<'loading' | 'authorized' | 'error'>('loading');
  const [popout, setPopout] = useState<ReactNode | null>(<ScreenSpinner />);

  useEffect(() => {
    async function fetchData() {
      try {
        const user = await bridge.send('VKWebAppGetUserInfo');
        setUser(user);

        const apiBaseURL = (import.meta.env.VITE_API_BASE_URL ?? '').trim();
        const authURL = apiBaseURL ? `${apiBaseURL}/auth/vk` : '/auth/vk';

        const response = await fetch(authURL, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ vk_id: user.id }),
        });
        if (!response.ok) {
          throw new Error(`Auth failed: ${response.status}`);
        }

        const auth = await response.json();
        if (auth.token) {
          localStorage.setItem('auth_token', auth.token);
          setAuthStatus('authorized');
        } else {
          throw new Error('JWT token is missing in response');
        }
      } catch (error) {
        setAuthStatus('error');
        console.error('VK auth error', error);
      } finally {
        setPopout(null);
      }
    }
    fetchData();
  }, []);

  return (
    <SplitLayout>
      <SplitCol>
        <View activePanel={activePanel}>
          <Home id="home" fetchedUser={fetchedUser} authStatus={authStatus} />
          <Surveys id="surveys" />
          <Persik id="persik" />
        </View>
      </SplitCol>
      {popout}
    </SplitLayout>
  );
};
