package ro.redeul.google.go.lang.psi.typing;

import ro.redeul.google.go.lang.psi.types.underlying.GoUnderlyingType;

public class GoTypePackage implements GoType {

    @Override
    public GoUnderlyingType getUnderlyingType() {
        return GoUnderlyingType.Undefined;
    }

    @Override
    public boolean isIdentical(GoType type) {
        return false;
    }

    @Override
    public void accept(Visitor visitor) {

    }
}
