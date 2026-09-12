import api from './api'

export function listUsers() {
  return api.get('/admin/users')
}

export function createUser(username, password, role) {
  return api.post('/admin/users', { username, password, role })
}

export function updateUser(userId, { newPassword, role }) {
  return api.patch(`/admin/users/${userId}`, { new_password: newPassword, role })
}