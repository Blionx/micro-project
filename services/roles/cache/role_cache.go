package cache

import (
	"errors"
	"sync"
	"roles/models"
)

var (
	roleCache []models.Role
	mu        sync.RWMutex
)

func SetRoles(roles []models.Role) {
	mu.Lock()
	defer mu.Unlock()
	roleCache = roles
}

func GetRoles() []models.Role {
	mu.RLock()
	defer mu.RUnlock()
	return roleCache
}

func FindById(id uint) (models.Role, error){
	var emptyRole models.Role
	for _, role := range roleCache {
		if(role.ID == id) {
			return role, nil
		}
	}
	return emptyRole, errors.New("Not Found")
}
