importScripts('/src/js/idb.js');
importScripts('/src/js/utility.js');

const CACHE_STATIC_NAME = 'static-v4';
const CACHE_DYNAMIC_NAME = 'dynamic-v2';
const STATIC_FILES = ['/',
'/index.html',
'/src/js/app.js',
'/src/js/idb.js',
'/src/js/feed.js',
'/src/js/utility.js',
'/src/js/material.min.js',
'/src/css/app.css',
'/src/css/feed.css',
'/src/images/main-image.jpg',
'https://fonts.googleapis.com/css?family=Roboto:400,700',
'https://fonts.googleapis.com/icon?family=Material+Icons',
'https://cdnjs.cloudflare.com/ajax/libs/material-design-lite/1.3.0/material.indigo-pink.min.css']

self.addEventListener('install', (event) => {
    console.log('[Service worker] Installing Service Worker...', event)
    event.waitUntil(
        caches.open(CACHE_STATIC_NAME).then((cache) => {
            console.log('[Cache] caching...');
            cache.addAll(STATIC_FILES);
        })
    )
})

self.addEventListener('activate', (event) => {
    console.log('[Service Worker] Activating Service Worker...', event);
    event.waitUntil(
        caches.keys().then(keyList => {
            return Promise.all(keyList.map(key => {
                if (key !== CACHE_STATIC_NAME && key !== CACHE_DYNAMIC_NAME) {
                    console.log('[Service Worker] Removing old cache');
                    return caches.delete(key);
                }
            }))
        })
    )
    return self.clients.claim();
})

function isInArray(string, array) {
    for (var i = 0; i < array.length; i++) {
        if (array[i] === string) {
            return true;
        }
    }
    return false;
}

self.addEventListener('fetch', (event) => {
    var url = 'https://pwagram-ff041-default-rtdb.firebaseio.com/posts';

    if (event.request.url.indexOf(url) > -1) {
        event.respondWith(fetch(event.request).then(res => {
            const clonedRes = res.clone();
            clearAllStorage('posts').then(() => {
                return clonedRes.json().then(data => Object.entries(data).forEach(([key, value]) => {
                    return writeData('posts', value);
                }))
            })
            return res;
        }))
    } else if (isInArray(event.request.url, STATIC_FILES)) {
        event.respondWith(caches.match(event.request));
    } else {
        event.respondWith(caches.match(event.request).then(response => {
            if (response) return response;
            return fetch(event.request);
        }))
    }
})