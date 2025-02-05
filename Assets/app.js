
class Router {
    constructor(routes) {
        this.routes = routes;
    }

    _getCurrentURL() {
        return window.location.pathname;
    }

    _loadInitialRoute() {
        const pathSegs = this._getCurrentURL().split('/').filter(seg => seg !== '');
        this.loadRoute(pathSegs);
    }

    _matchUrlToRoute(urlSegs) {
        const targetPath = '/' + urlSegs.join('/');
        return this.routes.find(route => route.path === targetPath);
    }

    loadRoute(urlSegs) {
        const matchedRoute = this._matchUrlToRoute(urlSegs);
        if (!matchedRoute) {
            throw new Error('Route not found');
        }
        // Assuming you have a DOM element with id 'app' to render views
        document.getElementById('app').innerHTML = matchedRoute.callback();
    }

    navigateTo(path) {
        window.history.pushState({}, '', path);
        const pathSegs = path.split('/').filter(seg => seg !== '');
        this.loadRoute(pathSegs);
    }
}

const routes = [
    { path: '/', callback: () => login() },
    { path: '/ragister', callback: () => register() },
];

const router = new Router(routes);

function login() {
    return `<div class="container">
            
            <div class="Paragraph">
                <h1>Talents Dashboard</h1>
                <h3>Lorem ipsum dolor sit amet consectetur adipisicing elit. Beatae cumque nulla recusandae eveniet quae impedit placeat consectetur autem alias? Dolore, accusantium quibusdam at atque libero voluptatum laboriosam numquam eum consequatur?</h3>
            </div>
            
            <div class="data-sender">
                <h1>LOGIN</h1>
                <div class="input-box">
                    <label for="Email Or Password">Email Or Username:</label>
                    <input type="text" placeholder="user or email" required>
                    <i class='bx bxs-user'></i>
                </div>

                <div class="input-box">
                    <label for="Password">Password:</label>
                    <input type="text" placeholder="Password" required>
                    <i class='bx bxs-lock-alt'></i>
                </div>
                
                <button id="login-btn">Login</button>
                
                <div class="register-link">
                    <p>Dont Have account ?
                        <a href="learn.zone01oujda.ma">Register</a>
                    </p>
                </div>

            </div>
            
    </div>`
}

function registerView() {
    return '<h1>About Page</h1>';
}

document.addEventListener('DOMContentLoaded', () => {
    router._loadInitialRoute();

    document.querySelectorAll('a').forEach(link => {
        link.addEventListener('click', function(e) {
            e.preventDefault();
            const path = this.getAttribute('href');
            router.navigateTo(path);
        });
    });
});

window.addEventListener('popstate', () => {
    router._loadInitialRoute();
});