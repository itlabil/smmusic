import api from './api'

export function listPlaylists() {
  return api.get('/playlists')
}

export function createPlaylist(name) {
  return api.post('/playlists', { name })
}

export function renamePlaylist(playlistId, name) {
  return api.patch(`/playlists/${playlistId}`, { name })
}

export function deletePlaylist(playlistId) {
  return api.delete(`/playlists/${playlistId}`)
}

export function listPlaylistSongs(playlistId) {
  return api.get(`/playlists/${playlistId}/songs`)
}

export function addSongToPlaylist(playlistId, songId) {
  return api.post(`/playlists/${playlistId}/songs`, { song_id: songId })
}

export function removeSongFromPlaylist(playlistId, songId) {
  return api.delete(`/playlists/${playlistId}/songs/${songId}`)
}

export function reorderPlaylist(playlistId, songIds) {
  return api.put(`/playlists/${playlistId}/reorder`, { song_ids: songIds })
}