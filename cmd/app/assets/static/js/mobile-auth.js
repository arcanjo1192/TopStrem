async function handleMobileLogin() {  
    try {  
        const response = await fetch('/auth/login');  
        const data = await response.json();  
          
        if (data.url) {  
            // Abre URL em Custom Tabs (implementação Android)  
            if (window.AndroidInterface) {  
                window.AndroidInterface.openCustomTab(data.url);  
            }  
        }  
    } catch (error) {  
        console.error('Erro ao obter URL de login:', error);  
        // Fallback para login normal  
        window.location.href = '/auth/login';  
    }  
}