export const handler = function(event, context) {
  try {
    // Verificar se é um evento válido do Cognito
    if (!event.request || !event.response) {
        throw new Error('Evento inválido do Cognito');
    }

    // Obter atributos do usuário
    const userAttributes = event.request.userAttributes;
    
    if (userAttributes) {
        const userRole = userAttributes['custom:role'] ?? 'user';    
        console.log(`Adicionando role '${userRole}' nas claims para usuário: ${userAttributes.email || userAttributes.sub}`);

        event.response = {
          "claimsAndScopeOverrideDetails": {
            "idTokenGeneration": {},
            "accessTokenGeneration": {
              "claimsToAddOrOverride": {
                'role': userRole, // Também adicionar sem o prefixo custom: para facilitar o uso
                'email': userAttributes.email
              },
              scopesToAdd: ["role"], // Não precisa
            }
          }
        };
        // Return to Amazon Cognito
        context.done(null, event);
    } else {
      throw "Empty User Attributes"      
    }
  } catch (error) {
    console.error('Erro no Pre Token Generation trigger:', error);

    // Em caso de erro, retornar o evento original para não bloquear a autenticação
    context.done(null, event);
  }
};