import React, { createContext, useContext } from 'react';
import type { Track } from '../types';
import usePlayer from '../hooks/usePlayer';

interface PlayerContextType {
  state: {
    currentTrack: Track | null;
    isPlaying: boolean;
    volume: number;
    progress: number;
    duration: number;
    playlist: Track[];
    currentTrackIndex: number;
  };
  togglePlay: () => void;
  playNext: () => void;
  playPrevious: () => void;
  setVolume: (volume: number) => void;
  setProgress: (progress: number) => void;
  setTrack: (track: Track, playlist?: Track[]) => void;
  currentTrack: Track | null;
  isPlaying: boolean;
  volume: number;
  progress: number;
}

const PlayerContext = createContext<PlayerContextType | undefined>(undefined);

export const PlayerProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const {
    state,
    play,
    pause,
    setTrack,
    setVolume,
    setProgress,
    playNext,
    playPrevious,
  } = usePlayer();

  const togglePlay = () => {
    if (state.isPlaying) {
      pause();
    } else {
      play();
    }
  };

  const value = {
    state,
    togglePlay,
    playNext,
    playPrevious,
    setVolume,
    setProgress,
    setTrack,
    currentTrack: state.currentTrack,
    isPlaying: state.isPlaying,
    volume: state.volume,
    progress: state.progress,
  };

  return <PlayerContext.Provider value={value}>{children}</PlayerContext.Provider>;
};

export const usePlayerContext = () => {
  const context = useContext(PlayerContext);
  if (context === undefined) {
    throw new Error('usePlayerContext must be used within a PlayerProvider');
  }
  return context;
};

export default PlayerContext; 