import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { login } from "../services/auth.api";
import { useAuth } from "../auth/useAuth";

export default function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const navigate = useNavigate();
  const { loginSuccess } = useAuth();

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();

    try {
      const data = await login(email, password);
      loginSuccess(data.accessToken, data.role);

      if (data.role === "admin") navigate("/admin");
      else if (data.role === "support") navigate("/support");
      else navigate("/customer");
    } catch {
      alert("Invalid credentials");
    }
  }

  return (
    <form onSubmit={handleSubmit}>
      <input value={email} onChange={e => setEmail(e.target.value)} />
      <input type="password" value={password} onChange={e => setPassword(e.target.value)} />
      <button>Login</button>
    </form>
  );
}
