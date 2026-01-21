export type Role = "admin" | "support" | "customer";

export type AuthContextType = {
  role: Role | null;
  loading: boolean;
  loginSuccess: (token: string, role: Role) => void;
  logout: () => void;
};

