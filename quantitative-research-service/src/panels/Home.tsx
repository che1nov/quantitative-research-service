import { FC } from 'react';
import {
  Panel,
  PanelHeader,
  Header,
  Button,
  Group,
  Cell,
  Div,
  Avatar,
  NavIdProps,
} from '@vkontakte/vkui';
import { UserInfo } from '@vkontakte/vk-bridge';
import { useRouteNavigator } from '@vkontakte/vk-mini-apps-router';

export interface HomeProps extends NavIdProps {
  fetchedUser?: UserInfo;
  authStatus: 'loading' | 'authorized' | 'error';
}

export const Home: FC<HomeProps> = ({ id, fetchedUser, authStatus }) => {
  const { photo_200, city, first_name, last_name } = { ...fetchedUser };
  const routeNavigator = useRouteNavigator();

  const statusLabel = {
    loading: 'Проверяем авторизацию...',
    authorized: 'Авторизация через VK выполнена',
    error: 'Ошибка авторизации. Проверьте backend и CORS.',
  }[authStatus];

  return (
    <Panel id={id}>
      <PanelHeader>Главная</PanelHeader>
      {fetchedUser && (
        <Group header={<Header size="s">User Data Fetched with VK Bridge</Header>}>
          <Cell before={photo_200 && <Avatar src={photo_200} />} subtitle={city?.title}>
            {`${first_name} ${last_name}`}
          </Cell>
        </Group>
      )}

      <Group header={<Header size="s">Статус интеграции</Header>}>
        <Cell>{statusLabel}</Cell>
      </Group>

      <Group header={<Header size="s">Действия</Header>}>
        <Div>
          <Button stretched size="l" mode="primary" onClick={() => routeNavigator.push('surveys')}>
            Создать опрос
          </Button>
        </Div>
        <Div>
          <Button stretched size="l" mode="secondary" onClick={() => routeNavigator.push('persik')}>
            Пройти демонстрационный экран
          </Button>
        </Div>
      </Group>
    </Panel>
  );
};
