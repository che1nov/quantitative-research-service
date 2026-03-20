const state = {
  token: "",
  survey: null,
  sessionID: "",
};

const tokenInput = document.getElementById("token");
const loadBtn = document.getElementById("loadBtn");
const saveBtn = document.getElementById("saveBtn");
const submitBtn = document.getElementById("submitBtn");
const statusEl = document.getElementById("status");
const surveySection = document.getElementById("surveySection");
const surveyTitle = document.getElementById("surveyTitle");
const surveyDescription = document.getElementById("surveyDescription");
const surveyForm = document.getElementById("surveyForm");

loadBtn.addEventListener("click", async () => {
  const token = tokenInput.value.trim();
  if (!token) {
    setStatus("Введите токен опроса");
    return;
  }

  try {
    setStatus("Загружаем опрос...");
    const surveyRes = await fetch(`/api/public/surveys/${encodeURIComponent(token)}`);
    if (!surveyRes.ok) {
      throw new Error("Опрос не найден");
    }
    state.survey = await surveyRes.json();
    state.token = token;

    const sessionRes = await fetch("/api/public/sessions", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ public_token: token }),
    });
    if (!sessionRes.ok) {
      throw new Error("Не удалось создать сессию");
    }

    const session = await sessionRes.json();
    state.sessionID = session.session_id;
    renderSurvey(state.survey);
    setStatus("Опрос загружен. Заполните ответы.");
  } catch (e) {
    setStatus(e.message || "Ошибка загрузки");
  }
});

saveBtn.addEventListener("click", async (e) => {
  e.preventDefault();
  if (!state.sessionID) return;

  const answers = collectAnswers();
  const res = await fetch("/api/public/sessions/progress", {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ session_id: state.sessionID, answers }),
  });
  setStatus(res.ok ? "Прогресс сохранён" : "Не удалось сохранить прогресс");
});

submitBtn.addEventListener("click", async (e) => {
  e.preventDefault();
  if (!state.sessionID) return;

  const answers = collectAnswers();
  const res = await fetch("/api/public/sessions/submit", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ session_id: state.sessionID, answers }),
  });
  setStatus(res.ok ? "Ответы отправлены. Спасибо!" : "Не удалось отправить ответы");
});

function setStatus(text) {
  statusEl.textContent = text;
}

function renderSurvey(survey) {
  surveySection.classList.remove("hidden");
  surveyTitle.textContent = survey.title;
  surveyDescription.textContent = survey.description || "";

  surveyForm.innerHTML = "";
  survey.questions.forEach((q) => {
    const box = document.createElement("div");
    box.className = "q";

    const title = document.createElement("h3");
    title.textContent = q.title;
    box.appendChild(title);

    if (q.type === "free_text") {
      const area = document.createElement("textarea");
      area.name = q.id;
      area.dataset.questionType = q.type;
      box.appendChild(area);
    } else {
      q.options.forEach((opt) => {
        const label = document.createElement("label");
        const input = document.createElement("input");
        input.type = q.type === "multi_choice" ? "checkbox" : "radio";
        input.name = q.id;
        input.value = opt.id;
        input.dataset.questionType = q.type;
        label.appendChild(input);
        label.appendChild(document.createTextNode(` ${opt.text || opt.id}`));
        box.appendChild(label);
      });
    }

    surveyForm.appendChild(box);
  });
}

function collectAnswers() {
  const answers = {};
  if (!state.survey) return answers;

  state.survey.questions.forEach((q) => {
    if (q.type === "free_text") {
      const el = surveyForm.querySelector(`textarea[name="${q.id}"]`);
      answers[q.id] = el && el.value ? [el.value] : [];
      return;
    }

    const checked = Array.from(surveyForm.querySelectorAll(`input[name="${q.id}"]:checked`));
    answers[q.id] = checked.map((item) => item.value);
  });

  return answers;
}
