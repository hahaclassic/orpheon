import { useState, useCallback, useRef, useEffect } from 'react';
import type { Track } from '../types';
import { apiService } from '../services/api';

interface PlayerState {
  currentTrack: Track | null;
  isPlaying: boolean;
  volume: number;
  progress: number;
  duration: number;
  playlist: Track[];
  currentTrackIndex: number;
}

const initialState: PlayerState = {
  currentTrack: null,
  isPlaying: false,
  volume: 1,
  progress: 0,
  duration: 0,
  playlist: [],
  currentTrackIndex: -1,
};

export const usePlayer = () => {
  const [state, setState] = useState<PlayerState>(initialState);
  const audioRef = useRef<HTMLAudioElement | null>(null);

  // Инициализация аудио элемента
  useEffect(() => {
    audioRef.current = new Audio();
    audioRef.current.preload = 'metadata';
    audioRef.current.volume = state.volume;

    const handleTimeUpdate = () => {
      if (audioRef.current) {
        setState(prev => ({
          ...prev,
          progress: audioRef.current?.currentTime || 0
        }));
      }
    };

    const handleLoadedMetadata = () => {
      if (audioRef.current) {
        setState(prev => ({
          ...prev,
          duration: audioRef.current?.duration || 0
        }));
      }
    };

    const handleEnded = () => {
      if (state.currentTrackIndex < state.playlist.length - 1) {
        playNext();
      } else {
        setState(prev => ({ ...prev, isPlaying: false, progress: 0 }));
      }
    };

    const handleError = (e: Event) => {
      console.error('Audio error:', e);
      setState(prev => ({ ...prev, isPlaying: false }));
    };

    audioRef.current.addEventListener('timeupdate', handleTimeUpdate);
    audioRef.current.addEventListener('loadedmetadata', handleLoadedMetadata);
    audioRef.current.addEventListener('ended', handleEnded);
    audioRef.current.addEventListener('error', handleError);

    return () => {
      if (audioRef.current) {
        audioRef.current.removeEventListener('timeupdate', handleTimeUpdate);
        audioRef.current.removeEventListener('loadedmetadata', handleLoadedMetadata);
        audioRef.current.removeEventListener('ended', handleEnded);
        audioRef.current.removeEventListener('error', handleError);
        audioRef.current.pause();
        audioRef.current.src = '';
      }
    };
  }, []);

  const play = useCallback(() => {
    if (audioRef.current) {
      const playPromise = audioRef.current.play();
      if (playPromise !== undefined) {
        playPromise
          .then(() => {
            setState(prev => ({ ...prev, isPlaying: true }));
          })
          .catch(error => {
            console.error('Error playing audio:', error);
            setState(prev => ({ ...prev, isPlaying: false }));
          });
      }
    }
  }, []);

  const pause = useCallback(() => {
    if (audioRef.current) {
      audioRef.current.pause();
      setState(prev => ({ ...prev, isPlaying: false }));
    }
  }, []);

  const setTrack = useCallback((track: Track, playlist: Track[] = []) => {
    if (audioRef.current) {
      const trackIndex = playlist.findIndex(t => t.id === track.id);
      const audioUrl = `http://localhost:8080/api/v1/tracks/${track.id}/audio`;
      
      // Сначала пауза текущего трека
      audioRef.current.pause();
      
      // Устанавливаем новый источник
      audioRef.current.src = audioUrl;
      
      // Добавляем заголовки для Range запроса
      audioRef.current.preload = 'metadata';
      
      // Устанавливаем тип контента
      audioRef.current.setAttribute('type', 'audio/mpeg');
      
      // Загружаем аудио
      audioRef.current.load();
      
      setState(prev => ({
        ...prev,
        currentTrack: track,
        isPlaying: false,
        playlist: playlist,
        currentTrackIndex: trackIndex,
        progress: 0
      }));

      // Автоматически начинаем воспроизведение после загрузки метаданных
      audioRef.current.addEventListener('loadedmetadata', () => {
        play();
      }, { once: true });
    }
  }, [play]);

  const playNext = useCallback(() => {
    if (state.currentTrackIndex < state.playlist.length - 1) {
      const nextTrack = state.playlist[state.currentTrackIndex + 1];
      setTrack(nextTrack, state.playlist);
    }
  }, [state.currentTrackIndex, state.playlist, setTrack]);

  const playPrevious = useCallback(() => {
    if (state.currentTrackIndex > 0) {
      const prevTrack = state.playlist[state.currentTrackIndex - 1];
      setTrack(prevTrack, state.playlist);
    }
  }, [state.currentTrackIndex, state.playlist, setTrack]);

  const setVolume = useCallback((volume: number) => {
    if (audioRef.current) {
      audioRef.current.volume = volume;
      setState(prev => ({ ...prev, volume }));
    }
  }, []);

  const setProgress = useCallback((progress: number) => {
    if (audioRef.current) {
      audioRef.current.currentTime = progress;
      setState(prev => ({ ...prev, progress }));
    }
  }, []);

  return {
    state,
    play,
    pause,
    setTrack,
    setVolume,
    setProgress,
    playNext,
    playPrevious,
    audioRef,
  };
};

export default usePlayer; 