export interface User {
  id: number;
  email: string;
  username: string;
  avatar?: string;
}

export interface Track {
  id: string;
  name: string;
  duration: number;
  track_number: number;
  coverUrl?: string;
}

export interface Playlist {
  id: number;
  name: string;
  coverImage?: string;
  trackCount: number;
  tracks: Track[];
}

export interface ApiResponse<T> {
  data: T;
  error?: string;
} 