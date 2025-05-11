import React, { createContext, useContext, useReducer, useCallback } from 'react';
import type { Track } from '../types';

interface PlayerState {
  currentTrack: Track | null;
  isPlaying: boolean;
  volume: number;
  progress: number;
}

interface PlayerContextType extends PlayerState {
  togglePlay: () => void;
  playNext: () => void;
  playPrevious: () => void;
  setVolume: (volume: number) => void;
  setProgress: (progress: number) => void;
  setTrack: (track: Track) => void;
}

type PlayerAction =
  | { type: 'TOGGLE_PLAY' }
  | { type: 'SET_VOLUME'; payload: number }
  | { type: 'SET_PROGRESS'; payload: number }
  | { type: 'SET_TRACK'; payload: Track };

const initialState: PlayerState = {
  currentTrack: null,
  isPlaying: false,
  volume: 1,
  progress: 0,
};

const PlayerContext = createContext<PlayerContextType | undefined>(undefined);

const playerReducer = (state: PlayerState, action: PlayerAction): PlayerState => {
  switch (action.type) {
    case 'TOGGLE_PLAY':
      return {
        ...state,
        isPlaying: !state.isPlaying,
      };
    case 'SET_VOLUME':
      return {
        ...state,
        volume: action.payload,
      };
    case 'SET_PROGRESS':
      return {
        ...state,
        progress: action.payload,
      };
    case 'SET_TRACK':
      return {
        ...state,
        currentTrack: action.payload,
        isPlaying: true,
        progress: 0,
      };
    default:
      return state;
  }
};

export const PlayerProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [state, dispatch] = useReducer(playerReducer, initialState);

  const togglePlay = useCallback(() => {
    dispatch({ type: 'TOGGLE_PLAY' });
  }, []);

  const playNext = useCallback(() => {
    // TODO: Implement next track logic
  }, []);

  const playPrevious = useCallback(() => {
    // TODO: Implement previous track logic
  }, []);

  const setVolume = useCallback((volume: number) => {
    dispatch({ type: 'SET_VOLUME', payload: volume });
  }, []);

  const setProgress = useCallback((progress: number) => {
    dispatch({ type: 'SET_PROGRESS', payload: progress });
  }, []);

  const setTrack = useCallback((track: Track) => {
    dispatch({ type: 'SET_TRACK', payload: track });
  }, []);

  const value = {
    ...state,
    togglePlay,
    playNext,
    playPrevious,
    setVolume,
    setProgress,
    setTrack,
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