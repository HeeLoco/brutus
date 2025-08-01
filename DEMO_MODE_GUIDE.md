# Demo Mode Usage Guide

## 🚀 How to Use Demo Mode

The application now supports complete storage account management alongside resource groups. Here's how to navigate the demo mode:

## 🔑 Login Options

When you first visit [http://localhost:5173](http://localhost:5173), you'll see the login form with these options:

### API Modes:
1. **Demo Mode** - Uses mock data, no backend required
2. **Go Backend** - Connects to Go backend (Port 8081) 
3. **Python Backend** - Connects to Python backend (Port 8000)
4. **TypeScript Backend** - Connects to TypeScript backend (Port 3000)
5. **Azure Direct** - Real Azure authentication with MSAL

## 📱 Using the Application

### After Login:
1. **Resource Groups**: Default landing page showing Azure resource groups
2. **Storage Accounts**: Click "💾 Storage Accounts" to manage storage accounts
3. **Navigation**: Seamless switching between resource types
4. **Debug**: Access environment and configuration debug info

### Available Features:
- ✅ **List Resources**: View all resource groups and storage accounts
- ✅ **Create Resources**: Full creation forms with validation
- ✅ **Delete Resources**: With confirmation dialogs
- ✅ **Detailed Views**: Complete resource information display
- ✅ **Demo Data**: Rich test data for both resource types

## 🚪 How to Exit Demo Mode

### Method 1: Logout Button (Recommended)
1. Click the **"Logout"** button in the top-right corner of any page
2. This will clear your session and return you to the login form
3. You can then select a different mode or change settings

### Method 2: Clear Login State
1. If logout doesn't work, go to [http://localhost:5173](http://localhost:5173)
2. If still showing the main app, open browser DevTools
3. Clear cookies/localStorage or use incognito mode

### Method 3: Browser Refresh
1. Simply refresh the page or close/reopen the browser tab
2. The application will return to the login form

## 🔧 Technical Details

### Demo Mode Features:
- **No Authentication Required**: Just enter any subscription ID format
- **Mock Data**: Pre-populated resource groups and storage accounts
- **Full CRUD Operations**: Create, read, update, delete (in memory only)
- **Realistic UI**: Same interface as production Azure management

### Storage Account Demo Data:
- 3 demo storage accounts across different resource groups
- Various configurations (Hot/Cool tiers, security settings)
- Complete endpoint information (blob, queue, table, file)
- Tags and metadata examples

## 🌐 Quick Test Workflow

1. **Login**: Select "Demo Mode" → Enter any subscription ID → Click "Connect"
2. **View Resource Groups**: See demo resource groups list
3. **Switch to Storage Accounts**: Click "💾 Storage Accounts" button  
4. **Explore Features**: Try creating, viewing, and deleting storage accounts
5. **Test Navigation**: Switch back to Resource Groups using navigation
6. **Exit**: Click "Logout" to return to login form

## 🎯 URLs for Testing

- **Frontend**: [http://localhost:5173](http://localhost:5173)
- **Backend Health**: [http://localhost:8080/health](http://localhost:8080/health)
- **Debug Page**: [http://localhost:5173/debug](http://localhost:5173/debug)

The demo mode logout issue has been fixed - you should now be able to easily exit demo mode using the logout button!