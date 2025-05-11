export function parseJwt(token: string): any {
  try {
    return JSON.parse(atob(token.split('.')[1]));
  } catch (e) {
    return null;
  }
}

export function isAdmin(token: string | null): boolean {
  if (!token) return false;
  const payload = parseJwt(token);
  return payload?.access_level === 2;
} 