// ===============================
// todo-app front-end (Vanilla JS)
//
// 役割：
// - クリック/変更イベントを拾う
// - fetchでサーバーへPOST(削除/完了切替/編集)
// - 成功したらDOMを更新
//
// ※ DB更新・期限切れ判定の最終責任は Go 側
// ===============================

console.log("script.js loaded");

/**
 * application/x-www-form-urlencoded で POST する共通関数
 */
function postForm(url, data) {
  const body = new URLSearchParams();
  for (const [k, v] of Object.entries(data)) {
    body.set(k, String(v));
  }

  return fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/x-www-form-urlencoded" },
    body,
  });
}

/**
 * YYYY/MM/DD or YYYY-MM-DD が今日より過去か？
 */
function isExpired(dateStr) {
  const today = new Date();
  today.setHours(0, 0, 0, 0);

  const normalized = (dateStr || "").replaceAll("/", "-");
  const d = new Date(normalized);
  if (Number.isNaN(d.getTime())) return false;

  d.setHours(0, 0, 0, 0);
  return d < today;
}

/**
 * todo-item の expired 表示を、表示中の日付から同期する
 * - 完了/未完了に関わらず、期限切れなら赤にしたい方針
 */
function syncExpiredClass(item) {
  if (!item) return;

  const dateSpan = item.querySelector(".editable-date");
  const dateText = dateSpan ? dateSpan.textContent.trim() : "";

  if (dateText && isExpired(dateText)) {
    item.classList.add("expired");
  } else {
    item.classList.remove("expired");
  }
}

/**
 * テキスト編集用 input に差し替え
 */
function replaceWithTextInput(hostEl, oldText) {
  const input = document.createElement("input");
  input.type = "text";
  input.value = oldText;
  input.className = "edit-input";

  hostEl.textContent = "";
  hostEl.appendChild(input);

  input.focus();
  input.select();
  return input;
}

// ===============================
// 削除（🗑️）
// ===============================
document.addEventListener("click", (e) => {
  const btn = e.target.closest(".delete-btn");
  if (!btn) return;

  const id = btn.dataset.id;
  if (!id) return;

  postForm("/delete", { id })
    .then((res) => {
      if (!res.ok) throw new Error("delete failed");
      btn.closest(".todo-item")?.remove();
    })
    .catch(console.error);
});

// ===============================
// 完了 / 未完了 切り替え（☑）
// ===============================
document.addEventListener("change", (e) => {
  const checkbox = e.target.closest(".toggle-checkbox");
  if (!checkbox) return;

  const id = checkbox.dataset.id;
  if (!id) return;

  const completed = checkbox.checked ? 1 : 0;

  postForm("/toggle", { id, completed })
    .then((res) => {
      if (!res.ok) throw new Error("toggle failed");

      const item = checkbox.closest(".todo-item");
      if (!item) return;

      const activeList = document.querySelector(".todo-list.active");
      const doneList = document.querySelector(".todo-list.done");
      if (!activeList || !doneList) return;

      // 完了状態に応じてDOM移動
      (completed ? doneList : activeList).appendChild(item);

      // ✅ 期限切れ表示を再同期（完了済みでも赤の方針）
      syncExpiredClass(item);
    })
    .catch((err) => {
      console.error(err);
      checkbox.checked = !checkbox.checked;
    });
});

// ===============================
// タイトルのインライン編集
// ===============================
document.addEventListener("click", (e) => {
  const titleEl = e.target.closest(".editable-title");
  if (!titleEl) return;
  if (titleEl.querySelector("input")) return;

  const id = titleEl.dataset.id;
  if (!id) return;

  const oldText = titleEl.textContent.trim();
  const input = replaceWithTextInput(titleEl, oldText);

  const save = () => {
    const newText = input.value.trim();
    if (newText === "" || newText === oldText) {
      titleEl.textContent = oldText;
      return;
    }

    postForm("/update-title", { id, title: newText })
      .then((res) => {
        if (!res.ok) throw new Error("update-title failed");
        titleEl.textContent = newText;
      })
      .catch((err) => {
        console.error(err);
        titleEl.textContent = oldText;
      });
  };

  input.addEventListener("blur", save);
  input.addEventListener("keydown", (ev) => {
    if (ev.key === "Enter") input.blur();
    if (ev.key === "Escape") titleEl.textContent = oldText;
  });
});

// ===============================
// 期限日のインライン編集
// ===============================
document.addEventListener("click", (e) => {
  const dateSpan = e.target.closest(".editable-date");
  if (!dateSpan) return;
  if (dateSpan.querySelector("input")) return;

  const id = dateSpan.dataset.id;
  if (!id) return;

  const oldText = dateSpan.textContent.trim(); // YYYY/MM/DD

  const input = document.createElement("input");
  input.type = "date";
  input.className = "edit-input";
  input.value = oldText.replaceAll("/", "-");

  dateSpan.textContent = "";
  dateSpan.appendChild(input);
  input.focus();

  const restore = (text) => {
    dateSpan.textContent = text;
  };

  const save = () => {
    const newDate = input.value; // YYYY-MM-DD
    if (!newDate) {
      restore(oldText);
      return;
    }

    postForm("/update-date", { id, due: newDate })
      .then((res) => {
        if (!res.ok) throw new Error("update-date failed");

        const formatted = newDate.replaceAll("-", "/");
        restore(formatted);

        // ✅ 日付変更直後に期限切れ表示も同期
        const item = dateSpan.closest(".todo-item");
        syncExpiredClass(item);
      })
      .catch((err) => {
        console.error(err);
        restore(oldText);
      });
  };

  input.addEventListener("blur", save);
  input.addEventListener("keydown", (ev) => {
    if (ev.key === "Enter") input.blur();
    if (ev.key === "Escape") restore(oldText);
  });
});

// ===============================
// 初期表示の期限切れclassを念のため同期したい場合（任意）
// すでにGoで付けているなら不要だけど、ズレ防止に残してOK
// ===============================
document.querySelectorAll(".todo-item").forEach(syncExpiredClass);
