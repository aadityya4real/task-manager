const BASE_URL = "https://task-manager-78yn.onrender.com";

// 🔐 SIGNUP
function signup() {
  const username = document.getElementById("username").value.trim();
  const password = document.getElementById("password").value.trim();

  if (!username || !password) {
    alert("Enter username & password");
    return;
  }

  fetch(`${BASE_URL}/signup`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password })
  })
  .then(res => res.json())
  .then(data => {
    console.log("Signup:", data);
    alert("Signup successful ✅");
  })
  .catch(err => console.error("Signup error:", err));
}

// 🔐 LOGIN
function login() {
  const username = document.getElementById("username").value.trim();
  const password = document.getElementById("password").value.trim();

  if (!username || !password) {
    alert("Enter username & password");
    return;
  }

  fetch(`${BASE_URL}/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password })
  })
  .then(res => res.json())
  .then(data => {
    console.log("LOGIN RESPONSE:", data);

    // 🔥 IMPORTANT: adjust this based on your backend
    const token = data.token || data.access_token || data.jwt;

    if (!token) {
      alert("Login failed ❌");
      return;
    }

    localStorage.setItem("token", token);

    document.getElementById("auth").style.display = "none";
    document.getElementById("app").style.display = "block";

    loadTasks();
  })
  .catch(err => console.error("LOGIN ERROR:", err));
}

// 📥 LOAD TASKS
function loadTasks() {
  fetch(`${BASE_URL}/tasks`, {
    headers: {
      "Authorization": "Bearer " + localStorage.getItem("token")
    }
  })
  .then(res => res.json())
  .then(result => {
    console.log("TASKS:", result);

    const tasks = result.data || [];

    const list = document.getElementById("taskList");
    list.innerHTML = "";

    tasks.forEach(task => {
      const li = document.createElement("li");

      li.innerHTML = `
        <input type="checkbox" 
          ${task.done ? "checked" : ""} 
          onchange="toggleTask(${task.id}, \`${task.title}\`, this.checked)" />

        <span style="flex:1; text-align:left; ${task.done ? 'text-decoration: line-through;' : ''}">
          ${task.title}
        </span>

        <button onclick="deleteTask(${task.id})">❌</button>
      `;

      list.appendChild(li);
    });
  })
  .catch(err => console.error("LOAD ERROR:", err));
}

// ➕ ADD TASK
function addTask() {
  const input = document.getElementById("taskInput");
  const title = input.value.trim();

  if (!title) {
    alert("Task cannot be empty");
    return;
  }

  fetch(`${BASE_URL}/tasks`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      "Authorization": "Bearer " + localStorage.getItem("token")
    },
    body: JSON.stringify({ title })
  })
  .then(res => res.json())
  .then(data => {
    console.log("ADD:", data);
    input.value = "";
    loadTasks();
  })
  .catch(err => console.error("ADD ERROR:", err));
}

// ❌ DELETE TASK
function deleteTask(id) {
  fetch(`${BASE_URL}/tasks?id=${id}`, {
    method: "DELETE",
    headers: {
      "Authorization": "Bearer " + localStorage.getItem("token")
    }
  })
  .then(() => loadTasks())
  .catch(err => console.error("DELETE ERROR:", err));
}

// 🔄 TOGGLE TASK
function toggleTask(id, title, done) {
  fetch(`${BASE_URL}/tasks?id=${id}`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
      "Authorization": "Bearer " + localStorage.getItem("token")
    },
    body: JSON.stringify({ title, done })
  })
  .then(() => loadTasks())
  .catch(err => console.error("TOGGLE ERROR:", err));
}

// 🚪 LOGOUT
function logout() {
  localStorage.removeItem("token");
  location.reload();
}

// 🔄 AUTO LOGIN
window.onload = function () {
  const token = localStorage.getItem("token");

  if (token) {
    document.getElementById("auth").style.display = "none";
    document.getElementById("app").style.display = "block";
    loadTasks();
  }
};