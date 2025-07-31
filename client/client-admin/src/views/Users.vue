<template>
  <div class="admin-users">
    <h1 class="header-name">Управление пользователями</h1>
    
    <!-- Фильтры -->
    <div class="filters">
      <input 
        v-model="searchQuery" 
        placeholder="Поиск"
        @input="fetchUsers(1)"
      />
      
      <select v-model="pageSize" @change="fetchUsers(1)">
        <option value="10">10 на странице</option>
        <option value="20">20 на странице</option>
        <option value="50">50 на странице</option>
      </select>
    </div>
    
    <!-- Таблица -->
    <div class="users-table">
      <div class="table-header">
          <div @click="sortBy('id')">ID</div>
          <div @click="sortBy('username')">Имя</div>
          <div @click="sortBy('email')">Email</div>
          <div @click="sortBy('created_at')">Дата регистрации</div>
      </div>
          
        <div class="rows-wrapper">
          <div class="table-row" v-for="user in users" :key="user.id">
           <div>{{ user.id }}</div>
            <div>{{ user.username }}</div>
            <div>{{ user.email }}</div>
            <div>{{ formatDate(user.created_at) }}</div>
         </div>
        </div>

     </div>
    
    <!-- Пагинация -->
    <div class="pagination">
      <button 
        v-for="page in pages" 
        :key="page"
        :class="{ active: currentPage === page }"
        @click="fetchUsers(page)"
      >
        {{ page }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';

const users = ref([]);
const searchQuery = ref('');
const currentPage = ref(1);
const pageSize = ref(20);
const totalUsers = ref(0);
const sortField = ref('id');
const sortDirection = ref('desc');

const pages = computed(() => {
  const totalPages = Math.ceil(totalUsers.value / pageSize.value);
  return Array.from({ length: totalPages }, (_, i) => i + 1);
});

async function fetchUsers(page = 1) {
  currentPage.value = page;
  
  const params = new URLSearchParams({
    page: page,
    page_size: pageSize.value,
    search: searchQuery.value,
    sort: sortField.value,
    order: sortDirection.value
  });
  
  try {
    const response = await fetch(`http://localhost:8080/api/admin/users?${params}`);
    const data = await response.json();
    
    users.value = data.users;
    totalUsers.value = data.pagination.total;
  } catch (error) {
    console.error('Ошибка загрузки:', error);
  }
}

function sortBy(field) {
  if (sortField.value === field) {
    sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc';
  } else {
    sortField.value = field;
    sortDirection.value = 'asc';
  }
  fetchUsers(1);
}

function formatDate(dateString) {
  return new Date(dateString).toLocaleDateString();
}

onMounted(() => fetchUsers(1));
</script>
<style scoped>
.admin-users {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.filters {
  margin: 20px 0;
  display: flex;
  gap: 15px;
}



.users-table {
  width: 100%;
  border-collapse: separate;
  box-shadow: 3px 6px 10px rgba(0, 0, 0, 0.3);
  border-radius: 12px;
  background-color: #383838;
  overflow: hidden;
}

.table-header {
  display: flex;
  justify-content: space-between;
}

.rows-wrapper {
  box-shadow: inset 1px 3px 5px  rgba(0, 0, 0, .4);
}

.header-name {
  font-size: 32px;
  text-align: center;
  padding-bottom: 12px;
}

.table-header div {
  cursor: pointer;
  transition: all .3s ease-in-out;
  height: 100%;
  width: 25%;
  background-color: #181818;
  padding: 20px;
  text-align: center;
  text-wrap: nowrap;
}

.table-header div:hover {
  background-color: #535353;
}

.table-row {
  display: flex;
  justify-content: space-between;
}

.table-row div {
  width: 25%;
  text-align: center;
  padding: 10px 0;
  border-right: 2px solid rgb(20, 20, 20, .5);
  border-bottom: 2px solid rgb(20, 20, 20, .5);
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.table-row div:last-child {
  border-right: none;
}
 


.pagination {
  margin-top: 20px;
  display: flex;
  gap: 5px;
}

.pagination button {
  padding: 5px 10px;
  border: 1px solid #ddd;
  background: white;
  cursor: pointer;
  border-radius: 4px;
}

.pagination button.active {
  background: #f1f500;
  color: black;
  border-color: #f1f500;
  
}
</style>