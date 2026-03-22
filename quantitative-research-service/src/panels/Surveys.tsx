import { FC, useEffect, useMemo, useState } from 'react';
import {
  Button,
  Cell,
  Div,
  FormItem,
  Group,
  Header,
  Input,
  NavIdProps,
  Panel,
  PanelHeader,
  PanelHeaderBack,
  Textarea,
} from '@vkontakte/vkui';
import { useRouteNavigator } from '@vkontakte/vk-mini-apps-router';

type Survey = {
  id: string;
  title: string;
  description: string;
  public_link: string;
};

export const Surveys: FC<NavIdProps> = ({ id }) => {
  const routeNavigator = useRouteNavigator();
  const apiBaseURL = useMemo(() => (import.meta.env.VITE_API_BASE_URL ?? '').trim(), []);
  const authToken = localStorage.getItem('auth_token') ?? '';

  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [status, setStatus] = useState('');
  const [surveys, setSurveys] = useState<Survey[]>([]);
  const [loading, setLoading] = useState(false);

  const apiURL = (path: string) => (apiBaseURL ? `${apiBaseURL}${path}` : path);

  const cabinetHeaders = () => {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    };

    if (authToken) {
      headers.Authorization = `Bearer ${authToken}`;
      return headers;
    }

    // Фолбэк для локальной отладки без JWT.
    headers['X-Internal-User'] = 'analyst';
    headers['X-CSRF-Token'] = import.meta.env.VITE_CSRF_TOKEN || 'dev-csrf-token';
    return headers;
  };

  const loadSurveys = async () => {
    try {
      setLoading(true);
      const response = await fetch(apiURL('/api/cabinet/surveys'), {
        method: 'GET',
        headers: cabinetHeaders(),
      });
      if (!response.ok) {
        throw new Error(`Ошибка загрузки списка: ${response.status}`);
      }
      const data: Survey[] = await response.json();
      setSurveys(data);
    } catch (error) {
      setStatus(`Не удалось загрузить опросы: ${String(error)}`);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void loadSurveys();
  }, []);

  const createSurvey = async () => {
    if (!title.trim()) {
      setStatus('Введите название опроса');
      return;
    }

    try {
      setLoading(true);
      setStatus('Создаем опрос...');
      const response = await fetch(apiURL('/api/cabinet/surveys'), {
        method: 'POST',
        headers: cabinetHeaders(),
        body: JSON.stringify({
          title,
          description,
          questions: [
            {
              id: 'q1',
              title: 'Как вам новый интерфейс?',
              type: 'single_choice',
              options: [
                { id: 'like', text: 'Нравится' },
                { id: 'dislike', text: 'Не нравится' },
              ],
            },
            {
              id: 'q2',
              title: 'Комментарий',
              type: 'free_text',
              options: [],
            },
          ],
        }),
      });

      if (!response.ok) {
        throw new Error(`Ошибка создания: ${response.status}`);
      }

      const created: Survey = await response.json();
      setStatus(`Опрос создан: ${created.title}`);
      setTitle('');
      setDescription('');
      await loadSurveys();
    } catch (error) {
      setStatus(`Не удалось создать опрос: ${String(error)}`);
    } finally {
      setLoading(false);
    }
  };

  const copyLink = async (link: string) => {
    try {
      await navigator.clipboard.writeText(link);
      setStatus('Публичная ссылка скопирована');
    } catch {
      setStatus('Не удалось скопировать ссылку');
    }
  };

  return (
    <Panel id={id}>
      <PanelHeader before={<PanelHeaderBack onClick={() => routeNavigator.back()} />}>Опросы</PanelHeader>

      <Group header={<Header size="s">Создание опроса</Header>}>
        <FormItem top="Название">
          <Input value={title} onChange={(event) => setTitle(event.target.value)} placeholder="Например: UX Survey" />
        </FormItem>
        <FormItem top="Описание">
          <Textarea value={description} onChange={(event) => setDescription(event.target.value)} placeholder="Кратко о цели опроса" />
        </FormItem>
        <Div>
          <Button size="l" stretched onClick={() => void createSurvey()} loading={loading}>
            Создать опрос
          </Button>
        </Div>
      </Group>

      <Group header={<Header size="s">Список опросов</Header>}>
        {surveys.length === 0 && <Cell subtitle="Создайте первый опрос">Пока опросов нет</Cell>}
        {surveys.map((survey) => (
          <Cell
            key={survey.id}
            subtitle={survey.description || 'Без описания'}
            after={<Button size="s" mode="secondary" onClick={() => void copyLink(survey.public_link)}>Скопировать ссылку</Button>}
          >
            {survey.title}
          </Cell>
        ))}
      </Group>

      <Group header={<Header size="s">Статус</Header>}>
        <Cell>{status || 'Готово к работе'}</Cell>
      </Group>
    </Panel>
  );
};
