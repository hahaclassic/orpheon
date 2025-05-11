import { useState, useCallback, useRef } from 'react';
import type { Track } from '../../core/infrastructure/services/api';

interface PlayerState {
  currentTrack: Track | null;
  isPlaying: boolean;
  volume: number;
  progress: number;
  duration: number;
}

const initialState: PlayerState = {
  currentTrack: null,
  isPlaying: false,
  volume: 50,
  progress: 0,
  duration: 0,
};

export const usePlayer = () => {
  const [state, setState] = useState<PlayerState>(initialState);
  const audioRef = useRef<HTMLAudioElement | null>(null);

  const play = useCallback(() => {
    if (audioRef.current) {
      audioRef.current.play();
      setState((prev) => ({ ...prev, isPlaying: true }));
    }
  }, []);

  const pause = useCallback(() => {
    if (audioRef.current) {
      audioRef.current.pause();
      setState((prev) => ({ ...prev, isPlaying: false }));
    }
  }, []);

  const setTrack = useCallback((track: Track) => {
    if (audioRef.current) {
      audioRef.current.src = track.coverUrl; // Assuming coverUrl is the audio URL
      audioRef.current.load();
      setState((prev) => ({ ...prev, currentTrack: track, isPlaying: false }));
    }
  }, []);

  const setVolume = useCallback((volume: number) => {
    if (audioRef.current) {
      audioRef.current.volume = volume / 100;
      setState((prev) => ({ ...prev, volume }));
    }
  }, []);

  const setProgress = useCallback((progress: number) => {
    if (audioRef.current) {
      audioRef.current.currentTime = progress;
      setState((prev) => ({ ...prev, progress }));
    }
  }, []);

  const handleTimeUpdate = useCallback(() => {
    if (audioRef.current) {
      setState((prev) => ({
        ...prev,
        progress: audioRef.current?.currentTime || 0,
      }));
    }
  }, []);

  const handleLoadedMetadata = useCallback(() => {
    if (audioRef.current) {
      setState((prev) => ({
        ...prev,
        duration: audioRef.current?.duration || 0,
      }));
    }
  }, []);

  const handleEnded = useCallback(() => {
    setState((prev) => ({ ...prev, isPlaying: false, progress: 0 }));
  }, []);

  return {
    state,
    play,
    pause,
    setTrack,
    setVolume,
    setProgress,
    handleTimeUpdate,
    handleLoadedMetadata,
    handleEnded,
    audioRef,
  };
};

export default usePlayer; 