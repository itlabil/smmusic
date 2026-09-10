import api from './api'

export function listSongs() {
  return api.get('/songs')
}

export function uploadSong(file, onUploadProgress) {
  const formData = new FormData()
  formData.append('file', file)
  return api.post('/songs/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
    onUploadProgress,
  })
}

export function getStreamUrl(songId, quality = 'standard') {
  const token = localStorage.getItem('token')
  return `/api/songs/${songId}/stream?quality=${quality}&token=${token}`
}

export function likeSong(songId) {
  return api.post(`/songs/${songId}/like`)
}

export function unlikeSong(songId) {
  return api.delete(`/songs/${songId}/like`)
}

export function listLikedSongs() {
  return api.get('/songs/liked')
}

export function recordPlay(songId) {
  return api.post(`/songs/${songId}/play`)
}

export function listRecentlyPlayed() {
  return api.get('/songs/recently-played')
}