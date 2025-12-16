# Publishing GoFetch to NPM

## Prerequisites

1. NPM account: https://www.npmjs.com/signup
2. Go installed (for building WASM)
3. Node.js 16+ installed

## Steps to Publish

### 1. Build the Package

```bash
npm run build
```

This will:
- Compile Go to WebAssembly (`gofetch.wasm`)
- Copy `wasm_exec.js` runtime
- Generate JavaScript wrapper (`gofetch.js`)
- Generate TypeScript definitions (`gofetch.d.ts`)

### 2. Test Locally (Optional)

You can test the package locally before publishing:

```bash
# In the gofetch directory
npm link

# In another project
npm link gofetch
```

### 3. Login to NPM

```bash
npm login
```

Enter your NPM credentials.

### 4. Publish to NPM

```bash
npm publish
```

Or for a scoped package:

```bash
npm publish --access public
```

### 5. Verify Publication

Visit: https://www.npmjs.com/package/gofetch

## Using the Published Package

Once published, users can install it:

```bash
npm install gofetch
```

### Usage Example

```javascript
import gofetch from 'gofetch';

// Create a client
const client = await gofetch.newClient();
client.setBaseURL('https://api.example.com');
client.setHeader('Authorization', 'Bearer token');

// Make requests
const response = await client.get('/users/1');
console.log(response.data);
```

## Package Version Management

To publish updates:

```bash
# Patch version (1.0.0 -> 1.0.1)
npm version patch

# Minor version (1.0.0 -> 1.1.0)
npm version minor

# Major version (1.0.0 -> 2.0.0)
npm version major

# Then publish
npm publish
```

## Troubleshooting

### Build Fails

- Ensure Go is installed: `go version`
- Ensure wasm_exec.js is found at `/usr/local/go/misc/wasm/wasm_exec.js`

### Publish Fails

- Check if package name is available: https://www.npmjs.com/package/gofetch
- Verify you're logged in: `npm whoami`
- Check package.json version is unique

## Package Structure

```
gofetch/
├── dist/                    # Published to NPM
│   ├── gofetch.js          # ES Module wrapper
│   ├── gofetch.d.ts        # TypeScript definitions
│   ├── gofetch.wasm        # Compiled Go code
│   └── wasm_exec.js        # Go WASM runtime
├── package.json            # NPM metadata
└── scripts/
    └── build-npm.js        # Build script
```

## Files Included in NPM Package

The `.npmignore` file ensures only necessary files are published:
- `dist/` folder
- `README.md`
- `LICENSE`

Everything else is excluded to keep the package size reasonable.
