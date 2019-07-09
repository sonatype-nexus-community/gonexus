package nexusrm

// "bytes"
// "text/template"

// User encapsulates the information about a Repository Manager user
type User struct {
	ID, FirstName, LastName, Password string
}

const groovyCreateUser = `import org.sonatype.nexus.security.*
import org.sonatype.nexus.security.user.*
import org.sonatype.nexus.security.role.*

def userId = '{{.ID}}'
User user = new User(userId: userId, firstName: '{{.FirstName}}', lastName: '{{.LastName}}', source: UserManager.DEFAULT_SOURCE, emailAddress: 'testUser@example.com', status: UserStatus.active, roles: [new RoleIdentifier(UserManager.DEFAULT_SOURCE, Roles.ADMIN_ROLE_ID)])
UserManager users = container.lookup(UserManager.class.name)
users.addUser(user, '{{.Password}}')`
