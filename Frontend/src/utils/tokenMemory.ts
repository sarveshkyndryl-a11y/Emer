let accessToken: string | null = null;

export const tokenMemory = {
  get(): string | null {
    return accessToken;
  },
  set(token: string) {
    accessToken = token;
  },
  clear() {
    accessToken = null;
  },
};
