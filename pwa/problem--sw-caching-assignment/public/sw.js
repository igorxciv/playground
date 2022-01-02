const CACHE_STATIC_NAME = 'static-v5';
const CACHE_DYNAMIC_NAME = 'dynamic-v1';

self.addEventListener('install', (event) => {
    console.log('[Service Worker] Installing...');
    event.waitUntil(caches.open(CACHE_STATIC_NAME).then((cache) => {
        console.log('caching...');
        cache.addAll([
            '/',
            '/index.html',
            'https://fonts.googleapis.com/css?family=Roboto:400,700',
            'https://fonts.googleapis.com/icon?family=Material+Icons',
            'https://cdnjs.cloudflare.com/ajax/libs/material-design-lite/1.3.0/material.indigo-pink.min.css',
            '/src/css/app.css',
            '/src/css/main.css',
            '/src/js/main.js',
            '/src/js/material.min.js',

        ])
    }))
})

self.addEventListener('activate', (event) => {
    console.log('[Service Worker] Activating...');
    event.waitUntil(caches.keys().then(keysList => {
        return Promise.all(keysList.map(key => {
            if (key !== CACHE_STATIC_NAME && key !== CACHE_DYNAMIC_NAME) {
                console.log('[Service Worker] Deleting old cache ', key);
                return caches.delete(key);
            }
        }))
    }))
    return self.clients.claim();
});

self.addEventListener('fetch', (event) => {
    event.respondWith(caches.match(event.request).then(response => {
        if (response) {
            return response;
        } else {
            return fetch(event.request).then(res => {
                caches.open(CACHE_DYNAMIC_NAME).then(cache => {
                    cache.put(event.request.url, res.clone());
                    return res;
                })
            })
        }
    }))
})