export type Role = "admin" | "support" | "customer";

export type AuthContextType = {
  role: Role | null;
  loginSuccess: (token: string, role: Role) => void;
  logout: () => void;
};
