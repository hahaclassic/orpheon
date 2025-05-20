import React, { createContext, useContext, useState } from 'react';
import type { ReactNode } from 'react';

type ContentType = "track" | "album" | "playlist" | "artist";

interface SearchState {
  query: string;
  results: any[];
  contentType: ContentType;
  country: string;
  genre: string;
  loading: boolean;
  error: string | null;
}

interface SearchContextType {
  searchState: SearchState;
  setSearchState: (state: SearchState) => void;
  clearSearch: () => void;
}

const initialState: SearchState = {
  query: '',
  results: [],
  contentType: 'track',
  country: '',
  genre: '',
  loading: false,
  error: null,
};

const SearchContext = createContext<SearchContextType | undefined>(undefined);

export const SearchProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [searchState, setSearchState] = useState<SearchState>(initialState);

  const clearSearch = () => {
    setSearchState(initialState);
  };

  return (
    <SearchContext.Provider value={{ searchState, setSearchState, clearSearch }}>
      {children}
    </SearchContext.Provider>
  );
};

export const useSearch = () => {
  const context = useContext(SearchContext);
  if (context === undefined) {
    throw new Error('useSearch must be used within a SearchProvider');
  }
  return context;
}; 