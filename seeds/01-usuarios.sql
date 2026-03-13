-- ============================================
-- SEEDER: Usuarios de prueba para login
-- ============================================
-- Contraseña para todos los usuarios: password123
-- Hash bcrypt generado con cost factor 12
--
-- IMPORTANTE: Este seeder solo CREA usuarios si NO EXISTEN
-- No sobrescribirá usuarios existentes (preserva contraseñas
-- y configuraciones modificadas)
--
-- Roles disponibles:
--   - admin (CON 2FA)
--   - supervisor (CON 2FA)
--   - agente (CON 2FA)
--   - contratista (CON 2FA)
--   - entidad (SIN 2FA) <-- ideal para pruebas rápidas
-- ============================================

INSERT INTO tickets_usuarios (
    nombres,
    apellidos,
    tipo_documento_id,
    numero_documento,
    email,
    telefono,
    regional_id,
    departamento_id,
    municipio_id,
    password,
    rol,
    origen,
    activo,
    created_at,
    updated_at
) VALUES
-- ========== ADMIN ==========
(
    'Alvaro Andres',
    'Alvarez Rodriguez',
    (SELECT id FROM tipo_documentos WHERE nombre = 'Cédula de Ciudadanía' LIMIT 1),
    '1117529658',
    'aaar529658@gmail.com',
    '3014593654',
    (SELECT id FROM regionales WHERE identificador = 100 LIMIT 1),
    (SELECT id FROM departamentos WHERE codigo_dane = '11' LIMIT 1),
    (SELECT id FROM municipios WHERE codigo_dane = '11001' LIMIT 1),
    '$2a$12$Dfgh.Uq8Ig.8Yal3YMRR9eGm5JYEjbSpSyIaXu5xM6hJ8Kk0bBgua',
    'admin',
    'seeder',
    true,
    NOW(),
    NOW()
),
(
    'Juan Camilo',
    'Puentes',
    (SELECT id FROM tipo_documentos WHERE nombre = 'Cédula de Ciudadanía' LIMIT 1),
    '1234567890',
    'jcpsandoval94@gmail.com',
    '3134212476',
    (SELECT id FROM regionales WHERE identificador = 100 LIMIT 1),
    (SELECT id FROM departamentos WHERE codigo_dane = '11' LIMIT 1),
    (SELECT id FROM municipios WHERE codigo_dane = '11001' LIMIT 1),
    '$2a$12$Dfgh.Uq8Ig.8Yal3YMRR9eGm5JYEjbSpSyIaXu5xM6hJ8Kk0bBgua',
    'admin',
    'seeder',
    true,
    NOW(),
    NOW()
),
-- ========== SUPERVISOR ==========
(
    'María Fernanda',
    'Supervisor Prueba',
    (SELECT id FROM tipo_documentos WHERE nombre = 'Cédula de Ciudadanía' LIMIT 1),
    '2345678901',
    'supervisor@sgi.gov.co',
    '3002345678',
    (SELECT id FROM regionales WHERE identificador = 100 LIMIT 1),
    (SELECT id FROM departamentos WHERE codigo_dane = '5' LIMIT 1),
    (SELECT id FROM municipios WHERE codigo_dane = '5001' LIMIT 1),
    '$2a$12$Dfgh.Uq8Ig.8Yal3YMRR9eGm5JYEjbSpSyIaXu5xM6hJ8Kk0bBgua',
    'supervisor',
    'seeder',
    true,
    NOW(),
    NOW()
),

-- ========== AGENTE ==========
(
    'Pedro Luis',
    'Agente Prueba',
    (SELECT id FROM tipo_documentos WHERE nombre = 'Cédula de Ciudadanía' LIMIT 1),
    '3456789012',
    'agente@sgi.gov.co',
    '3003456789',
    (SELECT id FROM regionales WHERE identificador = 200 LIMIT 1),
    (SELECT id FROM departamentos WHERE codigo_dane = '76' LIMIT 1),
    (SELECT id FROM municipios WHERE codigo_dane = '76001' LIMIT 1),
    '$2a$12$Dfgh.Uq8Ig.8Yal3YMRR9eGm5JYEjbSpSyIaXu5xM6hJ8Kk0bBgua',
    'agente',
    'seeder',
    true,
    NOW(),
    NOW()
),

-- ========== CONTRATISTA ==========
(
    'Ana María',
    'Contratista Prueba',
    (SELECT id FROM tipo_documentos WHERE nombre = 'Cédula de Ciudadanía' LIMIT 1),
    '4567890123',
    'contratista@sgi.gov.co',
    '3004567890',
    (SELECT id FROM regionales WHERE identificador = 300 LIMIT 1),
    (SELECT id FROM departamentos WHERE codigo_dane = '68' LIMIT 1),
    (SELECT id FROM municipios WHERE codigo_dane = '68001' LIMIT 1),
    '$2a$12$Dfgh.Uq8Ig.8Yal3YMRR9eGm5JYEjbSpSyIaXu5xM6hJ8Kk0bBgua',
    'contratista',
    'seeder',
    true,
    NOW(),
    NOW()
),

-- ========== ENTIDAD (sin 2FA - ideal para pruebas) ==========
(
    'Carlos Alberto',
    'Entidad Prueba',
    NULL,
    '900123456',
    'entidad@sgi.gov.co',
    '3005678901',
    (SELECT id FROM regionales WHERE identificador = 100 LIMIT 1),
    (SELECT id FROM departamentos WHERE codigo_dane = '25' LIMIT 1),
    (SELECT id FROM municipios WHERE codigo_dane = '25175' LIMIT 1),
    '$2a$12$Dfgh.Uq8Ig.8Yal3YMRR9eGm5JYEjbSpSyIaXu5xM6hJ8Kk0bBgua',
    'entidad',
    'seeder',
    true,
    NOW(),
    NOW()
)
ON CONFLICT (email) DO NOTHING;

-- ============================================
-- INSTRUCCIONES DE USO
-- ============================================
-- 1. Ejecutar este script en PostgreSQL:
--    psql -U postgres -d nombre_db -f seeds/01-usuarios.sql
--
-- 2. Para probar el login:
--    Email: entidad@sgi.gov.co (sin 2FA)
--    Password: password123
--
-- 3. Para probar con 2FA:
--    Email: admin@sgi.gov.co
--    Password: password123
--    (Requiere configurar Google Authenticator)
-- ============================================
