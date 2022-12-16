<template>
  <Suspense>
    <template #default>
      <router-view />
    </template>
    <template #fallback>
      Loading...
    </template>
  </Suspense>
  
</template>
